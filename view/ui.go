package view

import (
	"fmt"
	"log/slog"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"wwtool/viewmodel"
)

func BuildUI(win fyne.Window, vm *viewmodel.AppViewModel) fyne.CanvasObject {

	return container.NewAppTabs(
		container.NewTabItem("主页", homeUI(win, vm)),
		container.NewTabItem("帮助", helpUI()),
		// container.NewTabItem("测试", TestUI()),
	)
}

func homeUI(win fyne.Window, vm *viewmodel.AppViewModel) fyne.CanvasObject {
	return container.NewVBox(
		GameSelectorUI(win, vm),
		widget.NewSeparator(),
		changeServerUI(win, vm),
		widget.NewSeparator(),
		getGachaLinkUI(win, vm),
	)
}

// 游戏路径选择器 UI
func GameSelectorUI(win fyne.Window, vm *viewmodel.AppViewModel) fyne.CanvasObject {
	// 创建一个选择框，显示游戏路径列表
	selectWidget := widget.NewSelect(
		vm.GamePathList,
		func(selected string) {
			vm.SetSelectedID(slices.Index(vm.GamePathList, selected))
			slog.Info("选中路径", "path", selected, "id", vm.SelectedID)
		},
	)
	if vm.SelectedID >= 0 && vm.SelectedID < len(vm.GamePathList) {
		selectWidget.SetSelectedIndex(vm.SelectedID)
	}

	addBtn := widget.NewButton("添加游戏路径", func() {
		go func() {
			folder, _ := ShowFolderOpen(win)
			if folder != "" && !slices.Contains(vm.GamePathList, folder) {
				vm.AddGamePath(folder)
				fyne.DoAndWait(func() { selectWidget.SetOptions(vm.GamePathList) })
			}
		}()
	})
	delBtn := widget.NewButton("删除选中路径", func() {
		selected := selectWidget.Selected
		if selected == "" {
			dialog.ShowInformation("提示", "请先选择一个路径", win)
			return
		}
		vm.RemoveGamePath(selectWidget.SelectedIndex())
		selectWidget.SetOptions(vm.GamePathList)
		selectWidget.SetSelectedIndex(0) // 清除选中状态
		selectWidget.Refresh()
	})

	runBtn := widget.NewButton("启动", func() {
		err := vm.RunSelected()
		if err != nil {
			dialog.ShowError(err, win)
		} else {
			dialog.ShowInformation("启动成功", "程序已启动", win)
		}
	})

	return container.NewVBox(
		container.NewHBox(addBtn, delBtn, runBtn),
		selectWidget,
	)
}

// changeServerUI 创建软链接服务器管理界面
// 提供添加、删除软链接的功能
// 软链接格式为 {name: path}，name 是服务器简称，path 是软链接路径
// 选中某个服务器后，可以删除或修改其软链接路径
func changeServerUI(win fyne.Window, vm *viewmodel.AppViewModel) fyne.CanvasObject {

	server := []fyne.CanvasObject{
		widget.NewLabelWithStyle("官服/b服切换",
			fyne.TextAlignLeading,
			fyne.TextStyle{
				Bold:   true,
				Italic: true,
			},
		),
	}
	change := func(serverName string) {
		err := vm.CreateLinkToServer(serverName)
		info := ""
		if err != nil {
			info = fmt.Sprintf("切换到 %s 失败, %s", serverName, err.Error())
		} else {
			info = "成功切换到" + serverName
		}
		dialog.ShowInformation("切服结果", info, win)
	}
	for serverName := range vm.LinkServers {
		item := widget.NewButton(fmt.Sprintf("点击切换到: %s", serverName), func() {
			change(serverName)
		})
		server = append(server, item)
	}

	return container.NewVBox(server...)
}

// 获取抽卡链接 UI
func getGachaLinkUI(win fyne.Window, vm *viewmodel.AppViewModel) fyne.CanvasObject {

	linkUI := widget.NewLabel("")
	linkUI.Wrapping = fyne.TextWrapBreak

	return container.NewVBox(
		container.NewHBox(
			widget.NewButton("点击获取抽卡链接", func() {
				link := vm.GetGachaLink()
				if link == "" {
					link = "未获取到链接，请先在游戏里打开抽卡记录"
				}
				linkUI.SetText(link)
			}),
			widget.NewLabel("成功获取到抽卡链接后会自动复制到剪切板"),
		),
		linkUI,
	)
}
