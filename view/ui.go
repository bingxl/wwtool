package view

import (
	"log/slog"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
	tmpLinkAddr := ""
	tmpLinkName := binding.NewString()
	selectedServer := ""

	rg := widget.NewRadioGroup([]string{}, func(selected string) {
		if selected == "" {
			return
		}
		selectedServer = selected
		slog.Info("选中服务器", "server", selectedServer)
	})
	updateRg := func() {
		for server := range vm.LinkServers {
			if !slices.Contains(rg.Options, server) {
				rg.Append(server)
			}
		}
		rg.Refresh()
	}
	updateRg()

	addServer := func() {
		tmpLinkAddr = ""
		tmpLinkName.Set("")

		dialog.ShowForm(
			"添加软链接",
			"添加",
			"取消",
			[]*widget.FormItem{
				widget.NewFormItem("简称", widget.NewEntryWithData(tmpLinkName)),
				widget.NewFormItem("软链接路径",
					widget.NewButton(
						"点击",
						func() {
							dialog.ShowFolderOpen(
								func(reader fyne.ListableURI, err error) {
									if err != nil {
										slog.Error("选择软链接路径失败", "error", err)
										return
									}
									if reader == nil {
										slog.Info("用户取消选择文件夹")
										return
									}
									tmpLinkAddr = reader.Path()
									slog.Info("选择的软链接路径", "path", tmpLinkAddr)
									if err != nil {
									}
								},
								win,
							)
						},
					),
				),
			},
			func(isConfirmed bool) {
				slog.Info("添加软链接确认", "isConfirmed", isConfirmed, "tmpLinkAddr", tmpLinkAddr)

				name, _ := tmpLinkName.Get()
				if isConfirmed && tmpLinkAddr != "" && name != "" {
					vm.SetLinkServer(name, tmpLinkAddr)
					slog.Info("添加软链接", "name", name, "path", tmpLinkAddr)
					updateRg()
				}
			},
			win,
		)
	}

	return container.NewVBox(
		container.NewHBox(
			widget.NewButton("添加diff路径", addServer),
			widget.NewButton("删除选中软链接", func() {
				if selectedServer != "" {
					vm.SetLinkServer(selectedServer, "")
				}
			}),
			widget.NewButton("创建软链接", func() {
				slog.Info("创建软链接结果", "err", vm.CreateLinkToServer(selectedServer))
			}),
		),
		rg,
	)
}

// 获取抽卡链接 UI
func getGachaLinkUI(win fyne.Window, vm *viewmodel.AppViewModel) fyne.CanvasObject {
	gachaLink := binding.NewString()

	return container.NewVBox(
		container.NewHBox(
			widget.NewButton("获取抽卡链接", func() {
				gachaLink.Set(vm.GetGachaLink())
			}),
			widget.NewLabel("成功获取到抽卡链接后会自动复制到剪切板"),
		),

		widget.NewEntryWithData(gachaLink),
	)
}
