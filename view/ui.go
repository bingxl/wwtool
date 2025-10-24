package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"wwtool/viewmodel"
)

func BuildUI(win fyne.Window, vm *viewmodel.AppViewModel) fyne.CanvasObject {

	return container.NewAppTabs(
		container.NewTabItem(T("主页"), HomeUI(win, vm)),
		container.NewTabItem(T("库街区"), KujiequUI(win, vm)),
		container.NewTabItem(T("工具"), ThirdToolSelectorUI(win, vm)),
		container.NewTabItem(T("帮助"), HelpUI()),
		// container.NewTabItem(T("测试"), TestUI()),
	)
}

func HomeUI(win fyne.Window, vm *viewmodel.AppViewModel) fyne.CanvasObject {
	return container.NewVBox(
		GameSelectorUI(win, vm),
		widget.NewSeparator(),
		changeServerUI(win, vm),
		widget.NewSeparator(),
		getGachaLinkUI(win, vm),
	)
}

// title
//
// 传入的title 会自动调用 T 翻译
func title(title string) fyne.CanvasObject {
	return widget.NewLabelWithStyle(
		T(title),
		fyne.TextAlignLeading,
		fyne.TextStyle{
			Bold:   true,
			Italic: true,
		},
	)
}

// changeServerUI 创建软链接服务器管理界面
// 提供添加、删除软链接的功能
// 软链接格式为 {name: path}，name 是服务器简称，path 是软链接路径
// 选中某个服务器后，可以删除或修改其软链接路径
func changeServerUI(win fyne.Window, vm *viewmodel.AppViewModel) fyne.CanvasObject {
	server := []fyne.CanvasObject{
		title("官服/b服切换"),
		layout.NewSpacer(),
	}
	change := func(serverName string) {
		err := vm.CreateLinkToServer(serverName)
		info := ""
		if err != nil {
			info = T("切换到 %s 失败, %s", serverName, err.Error())
		} else {
			info = T("成功切换到") + serverName
		}
		// dialog.ShowInformation(T("切服结果"), info, win)
		ShowInfoWithAutoClose(T("切服结果"), info, win)
	}
	for serverName := range vm.LinkServers {
		item := widget.NewButton(T("切换到%s", T(serverName)), func() {
			change(serverName)
		})
		item.Importance = widget.HighImportance
		server = append(server, item)
	}

	return container.NewGridWithColumns(2, server...)
}

// 获取抽卡链接 UI
func getGachaLinkUI(win fyne.Window, vm *viewmodel.AppViewModel) fyne.CanvasObject {

	return container.NewHBox(
		widget.NewButton(T("点击获取抽卡链接"), func() {
			link := vm.GetGachaLink()
			if link == "" {
				link = T("未获取到链接，请先在游戏里打开抽卡记录")
			} else {
				link = link + "\n" + T("链接已复制到剪切板")
			}
			// dialog.ShowInformation("info", link, win)
			ShowInfoWithAutoClose("info", link, win)
		}),
		widget.NewLabel(T("成功获取到抽卡链接后会自动复制到剪切板")),
	)

}
