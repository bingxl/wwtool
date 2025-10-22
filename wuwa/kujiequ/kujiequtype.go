package kujiequ

import (
	"fmt"
	"strings"
	"time"
)

type KujieQuResponse[T any] struct {
	// 响应码， 200 成功； 202 token 过期
	Code    int    `json:"code"`
	Msg     string `json:"msg"`     // 提示信息
	Success bool   `json:"success"` //token有效时才有
	Data    T      `json:"data"`    // 具体的返回数据
}

type Token struct {
	Token   string `json:"token"`
	Devcode string `json:"devcode"`
}

// widget 相关结构体
// widget 结构体中的数据结构
type WidgetItemData struct {
	Name             string `json:"name"`             // 名称 ex. 结晶波片 活跃度
	Img              string `json:"img"`              //图标链接
	Key              string `json:"key"`              // 单位
	RefreshTimeStamp int64  `json:"refreshTimeStamp"` // 波片回满时间戳
	ExpireTimeStamp  int64  `json:"expireTimeStamp"`  // 资源过期时间戳
	Value            string `json:"value"`            // 值
	Cur              int    `json:"cur"`              // cur 当前数量
	Total            int    `json:"total"`            // 最大数量
}

func (w WidgetItemData) String() (r string) {
	refreshTime := time.Unix(w.RefreshTimeStamp, 0).Format("2006-1-2 15")
	r = fmt.Sprintf(
		"%s:%d/%d 刷新时间: %s",
		w.Name, w.Cur, w.Total, refreshTime,
	)

	return
}

// widget 请求返回的 data 结构
type WidgetResponseData struct {
	// 站双与鸣潮共同字段
	GameId     int    `json:"gameId"`     // 游戏id
	UserId     int    `json:"userId"`     // 库街区id
	ServerTime int64  `json:"serverTime"` // 服务器时间戳
	ServerId   string `json:"serverId"`   // 服务器id
	ServerName string `json:"serverName"` // 服务器名称
	SignInUrl  string `json:"signInUrl"`  // 签到页面链接
	HasSignIn  bool   `json:"hasSignIn"`  // 是否已签到
	SignInTxt  string `json:"signInTxt"`  // 签到提示文本
	RoleId     string `json:"roleId"`     // 角色id
	RoleName   string `json:"roleName"`   // 角色名称

	// GameId = 3 即鸣潮widget 专有
	EnergyData      WidgetItemData   `json:"energyData"`      // 结晶波片数据
	StoreEnergyData WidgetItemData   `json:"storeEnergyData"` // 结晶单质
	LivenessData    WidgetItemData   `json:"livenessData"`    // 每日活跃数据
	TowerData       WidgetItemData   `json:"towerData"`       //深塔
	SlashTowerData  WidgetItemData   `json:"slashTowerData"`  //海虚
	WeeklyData      WidgetItemData   `json:"weeklyData"`      // 战歌重奏
	WeeklyRougeData WidgetItemData   `json:"weeklyRougeData"` // 千道门扉
	BattlePassData  []WidgetItemData `json:"battlePassData"`  // 其他经验数据

	// GameId = 2 即站双widget专有
	ActionData          WidgetItemData   `json:"actionData"`          //血清数据
	DormData            WidgetItemData   `json:"dormData"`            // 宿舍委托数据
	ActiveData          WidgetItemData   `json:"activeData"`          // 每日活跃数据
	BossData            []WidgetItemData `json:"bossData"`            // 副本挑战数据
	ActionRecoverSwitch any              `json:"actionRecoverSwitch"` // 角色昵称，不常用，使用RoleName 替代
}

func (w WidgetResponseData) String() (format string) {
	arrFormat := func(arr []WidgetItemData) string {
		arrs := make([]string, len(arr))
		for i, v := range arr {
			arrs[i] = v.String()
		}
		return strings.Join(arrs, "\t")
	}

	switch w.GameId {
	case 2:
		format = fmt.Sprintf(
			"战双 %s    %s    %s    %s",
			w.RoleName, w.ActionData, w.DormData, arrFormat(w.BossData),
		)
	case 3:
		format = fmt.Sprintf(
			"鸣潮 %s    %s    %s    %s    %s    %s    %s    %s    %s",
			w.RoleName, w.EnergyData, w.LivenessData, w.TowerData, w.StoreEnergyData,
			w.SlashTowerData, w.WeeklyRougeData, w.WeeklyData,
			arrFormat(w.BattlePassData),
		)
	}
	return
}

// ------start角色信息
// 角色信息
type RoleInfo struct {
	Token      string `json:"token"`      // 令牌
	GameId     int    `json:"gameId"`     // 游戏ID
	RoleId     string `json:"roleId"`     // 角色ID
	RoleName   string `json:"roleName"`   // 角色名
	UserId     string `json:"userId"`     // 用户ID
	ServerId   string `json:"serverId"`   // 服务器ID
	ServerName string `json:"serverName"` // 服务器名
}

// 验证角色信息是否有效
func (r RoleInfo) IsValid() bool {
	if r.Token == "" ||
		r.GameId == 0 ||
		r.RoleId == "" ||
		r.RoleName == "" ||
		r.UserId == "" ||
		r.ServerId == "" ||
		r.ServerName == "" {
		return false
	}
	return true
}

// -------end 角色信息
