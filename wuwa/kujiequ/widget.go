package kujiequ

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

// 库街区鸣潮 小卡片

func widgetRequestHeader(token string) (h map[string]string) {
	h = map[string]string{
		"token":              token,
		"source":             "android",
		"sec-ch-ua-mobile":   "?1",
		"sec-ch-ua-platform": `"Android"`,
		"user-agent":         "Mozilla/5.0 (Linux; Android 14; 23127PN0CC Build/UKQ1.230804.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/127.0.6533.2 Mobile Safari/537.36 Kuro/2.2.0 KuroGameBox/2.2.0",
		"content-type":       "application/x-www-form-urlencoded",
		"devcode":            "61.178.245.214, Mozilla/5.0 (Linux; Android 14; 23127PN0CC Build/UKQ1.230804.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/127.0.6533.2 Mobile Safari/537.36 Kuro/2.2.0 KuroGameBox/2.2.0",

		"origin":          "https://web-static.kurobbs.com",
		"accept-language": "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
	}
	return
}

// 获取 role 角色对应的小卡片， 通过 role.GameId 区分是站双还是鸣潮
func (k *KujieQu) GetWidget(role RoleInfo) (widgetData WidgetResponseData, err error) {
	apis := map[int]string{
		2: Apis["punishWidget"],
		3: Apis["wuwaWidget"],
	}
	api := apis[role.GameId]

	payload := url.Values{
		"gameId":   []string{fmt.Sprintf("%d", role.GameId)},
		"roleId":   []string{role.RoleId},
		"serverId": []string{role.ServerId},
		"type":     []string{"2"},
		"sizeType": []string{"1"},
	}

	slog.Debug("info", "role", role)

	req, err := http.NewRequest("POST", api, bytes.NewBufferString(payload.Encode()))

	if err != nil {
		errMsg := fmt.Sprintf("获取鸣潮小卡片失败 角色名:%s, err: %s", role.RoleName, err.Error())

		return widgetData, errors.New(errMsg)
	}
	headers := widgetRequestHeader(k.token)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	// slog.Info("headers", "headers", req.Header)
	// return

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var resBody KujieQuResponse[WidgetResponseData]

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	slog.Debug("get widget raw data: " + string(bodyBytes))
	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&resBody)
	if err != nil {
		return
	}
	// slog.Info("resBody", "resBody", resBody)
	widgetData = resBody.Data
	return
}

func (k *KujieQu) GetAllWidgets() (widgets []WidgetResponseData, err error) {
	roles, err := k.FindAllRoles()
	if err != nil {
		slog.Error("获取角色失败" + err.Error())
		return
	}
	for _, role := range roles {
		widget, err := k.GetWidget(role)
		if err != nil {
			slog.Error("获取角色的小组件失败", "roleName", role.RoleName, "err", err)
			continue
		}
		slog.Info("成功获取的小组件", "roleName", role.RoleName, "widget", widget)
		widgets = append(widgets, widget)
	}

	return
}
