package card

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

// var T = i18n.Sprintf
func T(s string) string {
	// 模拟国际化函数
	// 实际使用中可以替换为真正的国际化函数
	return s
}

// -------------记录保存的格式
// []CardInfo 从新到旧排序
type CardSaveType struct {
	Info map[string]string                `json:"info"`
	List map[string]map[string][]CardInfo `json:"list"` // List: {playId: {cardType: []cardInfo} }
}

// 字符串化
func (c CardSaveType) String() string {
	b, _ := json.MarshalIndent(c, "", "\t")
	return string(b)
}

// 从字符串解析
func (c *CardSaveType) Parse(data string) error {
	err := json.Unmarshal([]byte(data), c)
	if err != nil {
		return err
	}
	return nil
}

func (c *CardSaveType) SetCards(playId string, data map[string][]CardInfo) {
	if data == nil {
		slog.Debug("CardSaveType.SetCards: data is nil")
		return
	}
	if c.List == nil {
		c.List = map[string]map[string][]CardInfo{}
	}
	if c.List[playId] == nil {
		c.List[playId] = map[string][]CardInfo{}
	}
	c.List[playId] = data
}

func (c CardSaveType) Valid() bool {
	// if len(c.List) == 0 {
	// 	return false
	// }
	// return true
	return len(c.List) != 0
}

// 鸣潮卡池一条记录
type CardInfo struct {
	CardPoolType string      `json:"cardPoolType"` // 卡池类型
	ResourceId   json.Number `json:"resourceId"`   // 资源id
	QualityLevel json.Number `json:"qualityLevel"` //	资源品质
	ResourceType string      `json:"resourceType"` // 资源类型
	Name         string      `json:"name"`         //	名字
	Count        json.Number `json:"count"`        //	数量
	Time         string      `json:"time"`         //	抽取时间 format: 2019-04-19 14:27:27
}

// 鸣潮抽卡记录获取与处理
type CardPool struct {
	params              map[string]string // 抽卡链接中的参数
	typePool            map[string]string // 卡池类型
	cn_REQUEST_URL      string            // 国内请求地址
	oversea_REQUEST_URL string            // 海外请求地址
	REQUEST_URL         string            // 请求地址
	Country             string            // 区域
	Store               CardStore         // 存储加载抽卡记录
}

func (c *CardPool) Start(params map[string]string, store CardStore) {

	c.cn_REQUEST_URL = "https://gmserver-api.aki-game2.com/gacha/record/query"
	c.oversea_REQUEST_URL = "https://gmserver-api.aki-game2.net/gacha/record/"
	c.Store = store
	c.params = params
	c.typePool = map[string]string{
		"1": T("角色活动唤取"), // key 为传递到请求的参数, value	为卡池类型放置到存储文件中
		"2": T("武器活动唤取"),
		"3": T("角色常驻唤取"),
		"4": T("武器常驻唤取"),
		"5": T("新手唤取"),
		"6": T("新手自选唤取"),
		"7": T("新手自选唤取（感恩定向唤取）"),
	}
	c.Country_URL()
	cards, err := c.LoadData()
	if err != nil {
		slog.Error(err.Error())
		cards = map[string][]CardInfo{}
	}
	for k, v := range c.typePool {
		newCard, err := c.Query(c.params, k)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		card := c.Update(cards[v], newCard)
		cards[v] = card
	}
	c.SaveData(cards)

}
func (c *CardPool) SaveData(records map[string][]CardInfo) error {
	slog.Debug("保存抽卡数据")
	card := CardSaveType{}
	card.SetCards(c.params["player_id"], records)
	return c.Store.Save(card)
}
func (c *CardPool) LoadData() (map[string][]CardInfo, error) {
	slog.Debug("加载抽卡数据")
	card, err := c.Store.Load()
	if err != nil {
		return nil, err
	}
	return card.List[c.params["player_id"]], nil
}

// 依据player_id获取区域和请求地址
func (c *CardPool) Country_URL(params ...map[string]string) string {
	slog.Info("running cardPoolRequest.Country_URL func")
	param := c.params
	if len(params) != 0 {
		param = params[0]
	}
	playId := param["player_id"]
	countrySet := map[byte][]string{
		'1': {c.cn_REQUEST_URL, "国服"},
		'6': {c.oversea_REQUEST_URL, "Eu"},
		'7': {c.oversea_REQUEST_URL, "Asia"},
		'8': {c.oversea_REQUEST_URL, "HMT (HK, MO, TW)"},
		'9': {c.oversea_REQUEST_URL, "SEA"},
	}
	if v, ok := countrySet[playId[0]]; ok {
		c.REQUEST_URL = v[0]
		c.Country = v[1]
		return v[0]
	}

	return ""
}

// @TODO: 请求抽卡记录
func (c *CardPool) Query(params map[string]string, cardPoolType string) ([]CardInfo, error) {
	slog.Info("running cardPoolRequest.Query func")

	// 参数转换
	reqParam := map[string]string{
		"playerId":     params["player_id"],
		"languageCode": params["lang"],
		// "gachaId":  params["gacha_id"],
		"serverId": params["svr_id"],
		// "svrArea":  params["svr_area"],
		"recordId":     params["record_id"],
		"cardPoolId":   params["resources_id"],
		"cardPoolType": cardPoolType,
	}

	paramJson, _ := json.Marshal(reqParam)
	slog.Info("请求参数", "paramJson", string(paramJson))
	res, err := http.Post(c.REQUEST_URL, "application/json", strings.NewReader(string(paramJson)))

	if err != nil {
		slog.Warn("请求抽卡记录失败" + err.Error())
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		slog.Info("请求失败")
		return nil, fmt.Errorf("获取抽卡记录失败,http 状态码: %d", res.StatusCode)
	}

	// 接口返回的数据类型
	var result struct {
		Code int        `json:"code"` // 0 成功 1 失败
		Data []CardInfo `json:"data"`
		Msg  string     `json:"msg"` // 成功时 success
	}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		slog.Info("解析数据失败" + err.Error())
		return nil, err
	}
	if result.Code != 0 {
		slog.Info("获取抽卡记录失败" + result.Msg)
		return nil, errors.New(result.Msg)
	}
	return result.Data, nil

}

// 更新记录
func (c *CardPool) Update(oldCard, newCard []CardInfo) []CardInfo {
	if len(oldCard) == 0 {
		return newCard
	}
	if len(newCard) == 0 {
		return oldCard
	}
	cards := make([]CardInfo, 0)
	timeFormat := "2006-01-02 15:04:05"
	newTime, err := time.Parse(timeFormat, newCard[len(newCard)-1].Time) // 从新到旧0].Time)
	if err != nil {
		slog.Error("解析时间失败", "time", newCard[len(newCard)-1].Time, "err", err.Error())
		return nil
	}
	// 鸣潮服务器中获取到的数据是从新到旧
	for i := 0; i < len(oldCard); i++ {
		oldTime, err := time.Parse(timeFormat, oldCard[i].Time)
		if err != nil {
			slog.Error("解析时间失败", "time", oldCard[i].Time, "err", err.Error())
			return nil
		}
		// 旧记录从新到旧的第一个 在 新纪录的最早时间 之前时候，就把旧记录剩余部分加入到新记录中(连同找到的第一个数据)
		if oldTime.Before(newTime) {
			cards = append(newCard, oldCard[i:]...)
			break
		}
	}

	return cards
}

func NewCardPool(params url.Values) *CardPool {
	p := map[string]string{}

	for k, v := range params {
		p[k] = v[0]
	}
	store := &CardStoreToFile{
		FilePath: filepath.Join("data", "cardPool-"+p["player_id"]+".json"),
	}
	c := new(CardPool)
	c.Start(p, store)
	return c
}
