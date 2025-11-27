package view

import (
	"errors"
	"log/slog"
	"strings"
	"wwtool/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// 游戏路径选择器 UI
func GameSelectorUI(win fyne.Window, vm *viewmodel.AppViewModel) fyne.CanvasObject {
	// 创建一个选择框，显示游戏路径列表
	selectWidget := widget.NewSelect(
		vm.GetGameNames(),
		func(selected string) {
			vm.SetSelectedID(selected)
			slog.Info("选中路径", "path", selected, "id", vm.BindGameSelectedIndex)
		},
	)
	initSelected, err := vm.BindGameSelectedIndex.Get()
	if err == nil && initSelected >= 0 {
		selectWidget.SetSelectedIndex(initSelected)
	}

	vm.BindGameNames.AddListener(binding.NewDataListener(func() {
		games, err := vm.BindGameNames.Get()
		if err == nil {
			selectWidget.SetOptions(games)
		}
	}))

	viewBtn := widget.NewButton(T("查看"), func() {
		game, _ := vm.GetCurrentSelectedGame()
		multiple := "✅"
		if !game.HasMultipleServer {
			multiple = "❌"
		}

		items := []*widget.FormItem{
			widget.NewFormItem(T("目录"), widget.NewLabel(game.Path)),
			widget.NewFormItem(T("多服"), widget.NewLabel(multiple)),
		}

		dialog.ShowForm(
			"info",
			T("确认"),
			T("取消"),
			items,
			func(bool) {},
			win,
		)
	})
	viewBtn.Importance = widget.MediumImportance

	delBtn := widget.NewButton(T("删除"), func() {
		selected := selectWidget.Selected
		if selected == "" {
			// dialog.ShowInformation(T("提示"), T("请先选择一个目录"), win)
			ShowInfoWithAutoClose(T("提示"), T("请先选择一个目录"), win)
			return
		}
		vm.RemoveGamePath(selectWidget.Selected)
		selectWidget.SetOptions(vm.GetGameNames())
		selectWidget.SetSelectedIndex(0) // 清除选中状态
		selectWidget.Refresh()
	})
	delBtn.Importance = widget.DangerImportance

	runBtn := widget.NewButton(T("启动"), func() {
		err := vm.RunSelected()
		if err != nil {
			dialog.ShowError(err, win)
		} else {
			// dialog.ShowInformation(T("启动成功"), T("程序已启动"), win)
			ShowInfoWithAutoClose(T("启动成功"), T("程序已启动"), win)
		}
	})
	runBtn.Importance = widget.HighImportance

	return container.NewVBox(
		addGameUI(win, vm),
		widget.NewSeparator(),
		title("选择与操作游戏"),
		container.NewGridWithColumns(
			4,
			selectWidget,
			delBtn,
			viewBtn,
			runBtn,
		),
	)
}

// 添加游戏路径
//
// onAdd - func(name string) 添加游戏成功时的回调, name 添加的别名
func addGameUI(win fyne.Window, vm *viewmodel.AppViewModel) fyne.CanvasObject {

	hasString := func(errInfo string) func(string) error {
		errInfo = T(errInfo)
		return func(s string) error {
			if strings.Trim(s, " ") == "" {
				return errors.New(errInfo)
			}
			return nil
		}

	}
	nameEntry := widget.NewEntry()
	nameEntry.Validator = hasString("请输入别名")
	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder(T("点击右侧按钮选择游戏目录"))
	pathEntry.Validator = hasString("请选择游戏目录")

	pathButton := widget.NewButton(T("选择游戏目录"), func() {
		// 执行选择游戏路径功能
		go func() {
			folder, err := ShowFolderOpen(win)
			if err == nil && folder != "" {
				// binding 线程安全， 不需要fyne.Do
				pathEntry.SetText(folder)
			}
		}()
	})
	hasMultiple := widget.NewCheck("", func(checked bool) {

	})

	form := widget.NewForm(
		widget.NewFormItem(T("别名"), nameEntry),
		widget.NewFormItem(T("目录"), container.NewAdaptiveGrid(2, pathEntry, pathButton)),
		widget.NewFormItem(T("多服"), hasMultiple),
	)
	form.SubmitText = T("添加游戏")
	form.OnSubmit = func() {
		if err := form.Validate(); err != nil {
			slog.Info("form validate error", "err", err)
			return
		}

		vm.AddGamePath(nameEntry.Text, pathEntry.Text, hasMultiple.Checked)
		// onAdd(name)
		// 添加后重置表单状态
		nameEntry.SetText("")
		pathEntry.SetText("")
		hasMultiple.SetChecked(false)

	}

	return widget.NewAccordion(
		widget.NewAccordionItem(T("添加游戏目录"), form),
	)
}
