package view

import (
	"log/slog"
	"wwtool/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// 游戏路径选择器 UI
func ThirdToolSelectorUI(win fyne.Window, vm *viewmodel.AppViewModel) fyne.CanvasObject {
	// 创建一个选择框，显示游戏路径列表
	selectWidget := widget.NewSelect(
		vm.GetThirdToolsName(),
		func(selected string) {
			vm.SetSelectedID(selected)
			slog.Info("选中路径", "path", selected, "id", vm.SelectedID)
		},
	)
	updateSelect := func() {
		selectWidget.SetOptions(vm.GetThirdToolsName())
		selectWidget.Refresh()
	}

	delBtn := widget.NewButton(T("删除"), func() {
		selected := selectWidget.Selected
		if selected == "" {
			ShowInfoWithAutoClose(T("提示"), T("请先选择一个目录"), win)
			return
		}
		vm.RemoveThirdTools(selectWidget.Selected)
		updateSelect()
	})
	delBtn.Importance = widget.DangerImportance

	runBtn := widget.NewButton(T("启动"), func() {
		vm.RunExe(selectWidget.SelectedIndex())
		ShowInfoWithAutoClose(T("已发送启动命令"), T("已发送启动命令"), win)

	})
	runBtn.Importance = widget.HighImportance

	onAddGame := func(name string) {
		// 刷新select 组件，设置选中的game
		vm.AddThirdTools(name)
		updateSelect()
	}

	return container.NewVBox(
		addThirdToolUI(win, vm, onAddGame),
		widget.NewSeparator(),
		title("选择与操作游戏"),
		container.NewGridWithColumns(
			3,
			selectWidget,
			delBtn,
			runBtn,
		),
	)
}

// 添加游戏路径
//
// onAdd - func(name string) 添加游戏成功时的回调, name 添加的别名
func addThirdToolUI(win fyne.Window, vm *viewmodel.AppViewModel, onAdd func(name string)) fyne.CanvasObject {

	bindPath := binding.NewString()

	pathEntry := widget.NewEntryWithData(bindPath)
	pathEntry.SetPlaceHolder(T("点击右侧按钮选择可执行文件"))

	pathButton := widget.NewButton(T("选择可执行文件"), func() {
		// filter := storage.NewExtensionFileFilter([]string{".exe"})
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				exe := reader.URI()
				slog.Info("fyne URI info",
					"path", exe.Path(),
					"name", exe.Name(),
					// "URI", exe,
				)
				bindPath.Set(exe.Path())
			}
		}, win)
		// 执行选择游戏路径功能

	})

	form := widget.NewForm(

		widget.NewFormItem(T("路径"), container.NewAdaptiveGrid(2, pathEntry, pathButton)),
	)
	form.SubmitText = T("添加")
	form.OnSubmit = func() {
		if err := form.Validate(); err != nil {
			slog.Info("form validate error", "err", err)
			return
		}
		path, _ := bindPath.Get()

		onAdd(path)
		// 添加后重置表单状态
		bindPath.Set("")
	}

	return widget.NewAccordion(
		widget.NewAccordionItem(T("添加可执行文件"), form),
	)
}
