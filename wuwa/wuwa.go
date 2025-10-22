package wuwa

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path"
	"slices"
	"wuwa/kujiequ"
)

type Token = kujiequ.Token

type UserInfo struct {
	UserOpenID string  `json:"userOpenID"`
	Tokens     []Token `json:"tokens"`
	GuildID    string  `json:"guildID"`
	ChannelID  string  `json:"channelID"`
}

type WuWaHandler interface {
	AddToken(userOpenID string, token string) bool
	GetTokens(userOpenID string) []string
	RemoveToken(userOpenID string, token string) bool
	KujiequSign(userOpenID string) (string, error)
	SetDataPath(path string) error
	GetUserInfo(userOpenID string) *UserInfo
}
type Handler struct {
	DataPath string
}

func (h *Handler) GetHeaders() map[string]string {
	return map[string]string{
		"ip":          "192.168.100.133",
		"version":     "2.2.5",
		"versioncode": "2250",
		"distinct_id": "da0b5c51-4627-4ea6-8c54-7ddce5cb6c31",
		"model":       "Redmi K30",
		"user-agent":  "okhttp/3.11.0",
	}
}

func (h *Handler) SetDataPath(path string) error {

	if path == "" {
		path = "./data"
	}
	h.DataPath = path
	return createDirIfNotExists(path)
}
func (h *Handler) readFile(filename string) (*UserInfo, error) {
	file := path.Join(h.DataPath, filename)
	data, err := os.ReadFile(file)
	if err != nil {
		slog.Error("读取文件失败", "file", file, "error", err)
		return nil, err
	}

	var obj *UserInfo
	err = json.Unmarshal(data, &obj)
	return obj, err
}
func (h *Handler) writeFile(filename string, data *UserInfo) error {
	marshalData, err := json.Marshal(data)
	if err != nil {
		slog.Error("序列化失败", "error", err)
		return err
	}
	file := path.Join(h.DataPath, filename)
	err = os.WriteFile(file, marshalData, 0644)
	if err != nil {
		slog.Error("写入文件失败", "file", file, "error", err)
		return err
	}
	return err
}

func (h *Handler) AddToken(userOpenID string, token string, devcode, guildID, channelID string) bool {

	obj, err := h.readFile(userOpenID + ".json")
	if err != nil {
		slog.Error(err.Error())
		obj = new(UserInfo)
		obj.UserOpenID = userOpenID
		obj.Tokens = make([]Token, 0)
	}
	if !slices.ContainsFunc(obj.Tokens, func(t Token) bool {
		return t.Token == token
	}) {
		obj.Tokens = append(obj.Tokens, Token{Token: token, Devcode: devcode})
	}

	obj.GuildID = guildID
	obj.ChannelID = channelID

	err = h.writeFile(userOpenID+".json", obj)
	return err == nil
}

func (h *Handler) GetTokens(userOpenID string) []Token {
	slog.Debug("in wuwa handler get tokens", "userOpenID", userOpenID)
	obj, err := h.readFile(userOpenID + ".json")
	if err != nil {
		slog.Error(err.Error())
		return []Token{}
	}

	return obj.Tokens
}

func (h *Handler) RemoveToken(userOpenID string, token string) bool {
	obj, err := h.readFile(userOpenID + ".json")
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	for i, v := range obj.Tokens {
		if v.Token == token {
			obj.Tokens = append(obj.Tokens[:i], obj.Tokens[i+1:]...)
			err = h.writeFile(userOpenID+".json", obj)
			return err == nil
		}
	}

	return false
}
func (h *Handler) GetUserInfo(userOpenID string) *UserInfo {
	obj, err := h.readFile(userOpenID + ".json")
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	return obj
}

func (h *Handler) GetWidgets(userOpenID string) (widgets []kujiequ.WidgetResponseData, err error) {
	tokens := h.GetTokens(userOpenID)
	if len(tokens) == 0 {
		err = errors.New("userOpenID:" + userOpenID + "没有token")
		return
	}

	for _, token := range tokens {
		headers := h.GetHeaders()
		headers["devcode"] = token.Devcode
		k := kujiequ.NewKujieQu(token, headers)
		widgetsTmp, err := k.GetAllWidgets()
		if err != nil {
			slog.Error("Handler.GetWidgets Kujiequ.GetAllWidgets error" + err.Error())
			continue
		}
		widgets = append(widgets, widgetsTmp...)
	}

	return
}

func (h *Handler) KujiequSign(userOpenID string) (string, error) {
	tokens := h.GetTokens(userOpenID)
	if len(tokens) == 0 {
		return "没有获取到任何token", nil
	}
	headers := h.GetHeaders()
	result := kujiequ.StartSign(tokens, headers)
	return result, nil
}

// 库街区签到与鸣潮抽卡记录

func Main() string {
	// startCardRecodeRequest()
	config, err := GetConfig()
	if err != nil {
		slog.Info("获取config失败" + err.Error())
		return "获取config失败" + err.Error()
	}
	return kujiequ.StartSign(config.Tokens, config.KujiequHeaders)
	// card.StartCardRecodeRequest(config.WuwaGamePath[0])
}

type Config struct {
	Tokens         []Token           `json:"tokens"`         // 库街区token列表
	WuwaGamePath   []string          `json:"wuwaGamePath"`   // 鸣潮游戏路径
	KujiequHeaders map[string]string `json:"kujiequHeaders"` // 库街区请求头
}

// 读取配置文件
// 可选参数 args 配置文件的位置，如果没有则使用当前目录下的config.json
func GetConfig(args ...string) (Config, error) {
	configFile := "config.json"
	if len(args) != 0 {
		configFile = args[0]
	}

	var config Config

	data, err := os.ReadFile(configFile)
	if err != nil {
		slog.Error("读取配置文件失败" + err.Error())
		return config, err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		slog.Error("解析配置文件失败" + err.Error())
		return config, err
	}

	return config, nil
}
