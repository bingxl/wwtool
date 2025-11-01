package kujiequ

// 库街区相关功能

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"net/url"
	"strconv"
	"strings"
	"time"
	"wuwa/req"
)

// -------start api config
// 库街区 api 请求地址
var KujiequHost = "https://api.kurobbs.com"
var Apis = map[string]string{
	"host":         KujiequHost,                                 // 服务器地址
	"sign":         KujiequHost + "/encourage/signIn/v2",        // 游戏签到地址
	"findRoles":    KujiequHost + "/user/role/findRoleList",     // 获取角色列表地址
	"bbsSign":      KujiequHost + "/user/signIn",                // bbs签到
	"wuwaWidget":   KujiequHost + "/gamer/widget/game3/refresh", // 鸣潮小组件
	"punishWidget": KujiequHost + "/gamer/widget/game2/refresh", // 站双小组件

	// 鸣潮资源简报
	"wuwaResourcePeriod":  KujiequHost + "/aki/resource/period/list", // 鸣潮资源简报索引
	"wuwaResourceMonth":   KujiequHost + "/aki/resource/month",       // 鸣潮每月资源简报
	"wuwaResourceWeek":    KujiequHost + "/aki/resource/week",        // 鸣潮每周资源简报
	"wuwaResourceVersion": KujiequHost + "/aki/resource/version",     // 鸣潮版本资源简报
}

// -------end api config

// -------start 常量/变量定义
var (
	GAME_IDS   = []int{2, 3}
	GAME_NAMES = map[int]string{
		2: "战双",
		3: "鸣潮",
	}
)

// -------end 常量/变量定义

// 库街区
type KujieQu struct {
	token   string
	devcode string
	headers map[string]string
}

func defaultHeader() map[string]string {
	return map[string]string{
		"ip":          "192.168.100.133",
		"version":     "2.2.5",
		"versioncode": "2250",
		"distinct_id": "da0b5c51-4627-4ea6-8c54-7ddce5cb6c31",
		"model":       "Redmi K30",
		"user-agent":  "okhttp/3.11.0",
	}
}

func (k *KujieQu) SetHeaders(headers map[string]string) {
	if k.headers == nil {
		k.headers = map[string]string{}

	}
	maps.Insert(k.headers, maps.All(defaultHeader()))
	maps.Insert(k.headers, maps.All(headers))

}
func (k *KujieQu) SetToken(token string, devcode string) {
	k.token = token
	k.devcode = devcode
}

// 查找当前token下 gameId 绑定的角色
func (k *KujieQu) FindRole(gameId int) ([]RoleInfo, error) {
	url := Apis["findRoles"]
	// 设置 body
	headers := map[string]any{
		"token":           k.token,
		"osversion":       "Android",
		"devcode":         k.devcode,
		"countrycode":     "CN",
		"ip":              k.headers["ip"],
		"model":           k.headers["model"],
		"source":          "android",
		"lang":            "zh-Hans",
		"version":         k.headers["version"],
		"versioncode":     k.headers["versioncode"],
		"content-type":    "application/x-www-form-urlencoded; charset=utf-8",
		"accept-encoding": "gzip",
		"user-agent":      k.headers["user-agent"],
	}

	body, err := req.Post(url, headers, strings.NewReader(fmt.Sprintf("gameId=%d", gameId)))
	// req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("gameId=%d", gameId)))
	if err != nil {
		return nil, err
	}

	var resBody KujieQuResponse[[]RoleInfo]
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&resBody)
	if err != nil {
		return nil, err
	}

	if resBody.Success {
		var roles []RoleInfo
		for _, v := range resBody.Data {
			v.Token = k.token
			if v.IsValid() {
				roles = append(roles, v)
			}
		}
		return roles, nil
	}

	return nil, fmt.Errorf("请求失败, code: %d, msg: %s", resBody.Code, resBody.Msg)
}

// 查找所有角色信息
// 传入一个可选 bool 参数 表示是否强制刷新
func (k *KujieQu) FindAllRoles(args ...bool) ([]RoleInfo, error) {
	forceFresh := len(args) > 0 && args[0]

	if roles, ok := GetGlobalRoles(k.token); ok && !forceFresh {
		return roles, nil
	}

	var roles []RoleInfo
	for _, v := range GAME_IDS {
		roleInfos, err := k.FindRole(v)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		roles = append(roles, roleInfos...)
	}
	AddGlobalRoles(k.token, roles)
	return roles, nil
}

// 查找账号绑定的角色 通过 gameId 过滤 2 站双； 3 鸣潮
func (k *KujieQu) FilterRoles(gameId int) []RoleInfo {
	roles, _ := k.FindAllRoles()
	filterRoes := []RoleInfo{}
	for _, role := range roles {
		if role.GameId == gameId {
			filterRoes = append(filterRoes, role)
		}
	}
	return filterRoes
}

// 角色签到
func (k *KujieQu) Sign(role RoleInfo) string {
	if !role.IsValid() {
		return fmt.Sprintf("角色信息不完整: %v", role)
	}
	// slog.Info("开始签到----", "roleName", role.RoleName)
	reqUrl := Apis["sign"]
	// 设置请求头
	headers := map[string]any{
		"token":              k.token,
		"sec-ch-ua":          "\"Not)A;Brand\";v=\"99\", \"Android WebView\";v=\"127\", \"Chromium\";v=\"127\"",
		"source":             "android",
		"sec-ch-ua-mobile":   "?1",
		"user-agent":         "Mozilla/5.0 (Linux; Android 12; Redmi k30 Build/UKQ1.210908.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/129.0.6668.100 Mobile Safari/537.36 Kuro/2.2.5 KuroGameBox/2.2.5",
		"content-type":       "application/x-www-form-urlencoded",
		"devcode":            "183.17.51.208, Mozilla/5.0 (Linux; Android 12; Redmi k30 Build/UKQ1.210908.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/129.0.6668.100 Mobile Safari/537.36 Kuro/2.2.5 KuroGameBox/2.2.5",
		"sec-ch-ua-platform": "\"Android\"",
		"origin":             "https://web-static.kurobbs.com",
		"x-requested-with":   "com.kurogame.kjq",
		"accept-language":    "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
	}
	// 设置 body
	payload := url.Values{
		"roleId":   {role.RoleId},
		"serverId": {role.ServerId},
		"gameId":   {strconv.Itoa(role.GameId)},
		"userId":   {role.UserId},
		"reqMonth": {time.Now().Format("01")},
	}
	// slog.Info("in sign payload", "payload", payload.Encode())
	body, err := req.Post(reqUrl, headers, strings.NewReader(payload.Encode()))
	if err != nil {
		return err.Error()
	}

	var resBody KujieQuResponse[any]
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&resBody)
	if err != nil {
		return "解析响应失败" + err.Error()
	}
	slog.Debug(fmt.Sprintf("response: %+v", resBody))
	return fmt.Sprintf("%v", resBody.Msg)
}

// 库街区社区签到
func (k *KujieQu) BbsSign() string {
	api := Apis["bbsSign"]
	headers := map[string]any{
		"token":        k.token,
		"osversion":    "Android",
		"devcode":      k.devcode,
		"distinct_id":  k.headers["distinct_id"],
		"countrycode":  "CN",
		"ip":           k.headers["ip"],
		"model":        k.headers["model"],
		"source":       "android",
		"lang":         "zh-Hans",
		"version":      k.headers["version"],
		"versioncode":  k.headers["versioncode"],
		"content-type": "application/x-www-form-urlencoded",
		"user-agent":   k.headers["user-agent"],
	}

	body, err := req.Post(api, headers, strings.NewReader("gameId=2"))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err).Error()
	}
	var resBody KujieQuResponse[any]
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&resBody)
	if err != nil {
		return "解析响应失败" + err.Error()
	}
	return resBody.Msg
}

func NewKujieQu(token Token, headers map[string]string) KujieQu {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["devcode"] = token.Devcode

	k := KujieQu{token: token.Token, devcode: token.Devcode}
	k.SetHeaders(headers)
	return k
}

// 开始获取所有角色信息并签到
func StartSign(tokens []Token, headers map[string]string) string {
	if len(tokens) == 0 {
		return "没有获取到任何token"
	}

	ret := ""
	for _, token := range tokens {
		k := NewKujieQu(token, headers)
		roles, err := k.FindAllRoles()
		if err != nil {
			slog.Error(err.Error())
			return "获取角色信息失败" + err.Error()
		}
		slog.Debug("获取角色信息成功", "roles", roles)
		for _, role := range roles {
			slog.Info("-------开始签到", "gameId", role.GameId, "roleName", role.RoleName)
			r := k.Sign(role)
			info := fmt.Sprintf("%v:%v => %v", GAME_NAMES[role.GameId], role.RoleName, r)
			slog.Info(r)
			ret += info + "\n"
			// 每个角色签到间隔1毫秒，防止请求过于频繁
			time.Sleep(time.Millisecond)
		}
		slog.Info("-------开始库街区签到")
		bbsResult := k.BbsSign()
		slog.Info(bbsResult)
		ret += "库街区签到结果:" + bbsResult + "\n"

		// 每个token之间间隔，防止请求过于频繁
		time.Sleep(time.Millisecond * 1)
	}
	return ret

}
