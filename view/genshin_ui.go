package view

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"wwtool/i18n"
	"wwtool/viewmodel"
)

// GenshinUI 创建原神游戏管理界面
func GenshinUI(win fyne.Window) fyne.CanvasObject {
	vm := viewmodel.NewGenshinViewModel()
	// 游戏文件路径显示

	// 选择游戏文件按钮
	selectFileBtn := widget.NewButton(i18n.T("点击选择原神游戏本体文件"), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				path := reader.URI().Path()
				vm.SetGenshinPath(path)

			}
		}, win)
	})

	// 服务器切换按钮
	changeServerBtn := func(serverName byte, serverDisplayName string) *widget.Button {
		btn := widget.NewButton(i18n.T("切换到%s", serverDisplayName), func() {
			if v, _ := vm.GenshinPath.Get(); v == "" {
				dialog.ShowError(errors.New(i18n.T("请先选择游戏文件路径")), win)
				return
			}

			vm.SwitchServer(serverName)

			dialog.ShowInformation(
				i18n.T("切换成功"),
				i18n.T("已切换到%s服务器", serverDisplayName),
				win,
			)
		})
		btn.Importance = widget.HighImportance
		return btn
	}

	officialBtn := changeServerBtn('g', i18n.T("官服"))
	bilibiliBtn := changeServerBtn('b', i18n.T("B服"))

	// 布局
	fileSelectContainer := container.NewVBox(
		title(i18n.T("原神游戏设置")),

		selectFileBtn,
		widget.NewLabelWithData(vm.GenshinPath),
		widget.NewSeparator(),
	)

	serverContainer := container.NewVBox(
		title(i18n.T("服务器切换")),
		container.NewHBox(officialBtn, bilibiliBtn),
		widget.NewSeparator(),
	)

	return container.NewVBox(
		fileSelectContainer,
		serverContainer,
	)
}
