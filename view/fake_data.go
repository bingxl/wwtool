package view

import (
	"fmt"
	"wuwa/kujiequ"
)

func fakeWidgetItem(name string, cur, total int) kujiequ.WidgetItemData {
	return kujiequ.WidgetItemData{
		Name:             name,
		Cur:              cur,
		Total:            total,
		RefreshTimeStamp: 1761010701,
	}
}

var kujiequWidgetsFakeData = []kujiequ.WidgetResponseData{
	kujiequ.WidgetResponseData{
		GameId:     3,
		UserId:     10000000,
		ServerTime: 1760948124,
		ServerId:   "76402e5b20be2c39f095a152090afddc",
		ServerName: "鸣潮",
		SignInTxt:  "已完成签到",
		HasSignIn:  true,
		RoleId:     "125020135",
		RoleName:   "可乐",
		EnergyData: fakeWidgetItem("结晶波片", 66, 240),

		LivenessData: fakeWidgetItem("活跃度", 100, 100),
		BattlePassData: []kujiequ.WidgetItemData{
			fakeWidgetItem("电台等级", 20, 0),
			fakeWidgetItem("本周经验", 3000, 12000),
		},
		StoreEnergyData: fakeWidgetItem("结晶单质", 345, 480),
		TowerData:       fakeWidgetItem("逆境深塔·深境区", 12, 36),
		SlashTowerData:  fakeWidgetItem("冥歌海墟·禁忌海域", 0, 0),
		WeeklyData:      fakeWidgetItem("战歌重奏", 0, 3),
		WeeklyRougeData: fakeWidgetItem("千道门扉的异想", 6000, 6000),
	},
	{
		GameId:     3,
		UserId:     10000000,
		ServerTime: 1760948124,
		ServerId:   "76402e5b20be2c39f095a152090afddc",
		ServerName: "鸣潮",
		SignInTxt:  "已完成签到",
		HasSignIn:  true,
		RoleId:     "125020135",
		RoleName:   "bing",
		EnergyData: fakeWidgetItem("结晶波片", 66, 240),

		LivenessData: fakeWidgetItem("活跃度", 100, 100),
		BattlePassData: []kujiequ.WidgetItemData{
			fakeWidgetItem("电台等级", 20, 0),
			fakeWidgetItem("本周经验", 3000, 12000),
		},
		StoreEnergyData: fakeWidgetItem("结晶单质", 345, 480),
		TowerData:       fakeWidgetItem("逆境深塔·深境区", 12, 36),
		SlashTowerData:  fakeWidgetItem("冥歌海墟·禁忌海域", 0, 0),
		WeeklyData:      fakeWidgetItem("战歌重奏", 0, 3),
		WeeklyRougeData: fakeWidgetItem("千道门扉的异想", 6000, 6000),
	},
}

func kujiequWidgetToMD(widget kujiequ.WidgetResponseData) string {
	md := fmt.Sprintf(
		`%s %s
是否已签到:%v
%s
%s
%s
%s
%s
%s
%s
`,
		widget.ServerName, widget.RoleName,
		widget.HasSignIn, widget.LivenessData, widget.EnergyData,
		widget.StoreEnergyData, widget.TowerData, widget.SlashTowerData, widget.WeeklyData,
		widget.WeeklyRougeData,
	)

	return md
}
