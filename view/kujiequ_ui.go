package view

import (
	"fmt"
	"log/slog"
	"time"
	"wuwa/kujiequ"
	"wwtool/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func KujiequUI(win fyne.Window, vm *viewmodel.AppViewModel) fyne.CanvasObject {
	// for _, token := range tokens {
	// 	widgets := vm.GetKujiequWidgets(token)
	// }

	// 创建一个选择框，显示游戏路径列表
	selectWidget := widget.NewSelect(
		vm.GetTokens(),
		func(selected string) {
			vm.SelectedKujiequToken.Set(selected)
			slog.Info("选中的token", "token", selected)
		},
	)
	updateSelect := func() {
		selectWidget.SetOptions(vm.GetTokens())
		selectWidget.Refresh()
	}

	delBtn := widget.NewButton(T("删除"), func() {
		selected := selectWidget.Selected
		if selected == "" {
			ShowInfoWithAutoClose(T("提示"), T("请先选择一个token"), win)
			return
		}
		vm.DeleteToken(selectWidget.Selected)
		updateSelect()
	})
	delBtn.Importance = widget.DangerImportance

	runBtn := widget.NewButton(T("获取小组件"), func() {
		// vm.RunExe(selectWidget.SelectedIndex())
		// ShowInfoWithAutoClose(T("已发送启动命令"), T("已发送启动命令"), win)
		if selectWidget.Selected != "" {
			slog.Info("点击了获取小组件，并且token 不为空")
			now := time.Now().Unix()
			// 使用当前时间作为参数来刷新 widgets
			vm.RefreshKujiequWidget.Set(fmt.Sprintf("%d", now))

			// wids := vm.GetKujiequWidgets(selectWidget.Selected)
			// vm.SelectedKujiequToken.
		}

	})
	runBtn.Importance = widget.HighImportance

	onAddToken := func(token string, devcode string) {
		// 刷新select 组件，设置选中的game
		vm.AddToken(token, devcode)
		updateSelect()
	}

	return container.NewVBox(
		addTokenUI(win, vm, onAddToken),
		widget.NewSeparator(),
		title("选择token"),
		container.NewGridWithColumns(
			3,
			selectWidget,
			delBtn,
			runBtn,
		),
		container.NewStack((widgetGridObj(vm))),
	)
}

// 渲染库街区小组件
func widgetGridObj(vm *viewmodel.AppViewModel) fyne.CanvasObject {
	kujiequWidgetsUI := container.NewGridWithColumns(2)
	kujiequLayout := widget.NewAccordion(
		widget.NewAccordionItem("小组件", kujiequWidgetsUI),
	)
	kujiequLayout.OpenAll()

	nameAlias := func(name string) string {
		alias := map[string]string{
			"逆境深塔·深境区":  "深塔",
			"冥歌海墟·禁忌海域": "海墟",
			"千道门扉的异想":   "千道门扉",
		}
		if aliasName, ok := alias[name]; ok {
			return aliasName
		}
		return name
	}
	itemRender := func(items ...kujiequ.WidgetItemData) []fyne.CanvasObject {
		renders := make([]fyne.CanvasObject, len(items))
		for index, item := range items {
			renders[index] = widget.NewLabel(
				fmt.Sprintf("%s: %d/%d", nameAlias(item.Name), item.Cur, item.Total),
			)
		}

		return renders
	}
	//
	card := func(resp kujiequ.WidgetResponseData) fyne.CanvasObject {
		var renderItems = []kujiequ.WidgetItemData{}
		switch resp.GameId {
		// 站双
		case 2:
			renderItems = []kujiequ.WidgetItemData{
				resp.ActionData, resp.DormData, resp.ActionData,
			}
			renderItems = append(renderItems, resp.BossData...)
		// 鸣潮
		case 3:
			renderItems = []kujiequ.WidgetItemData{
				resp.EnergyData, resp.LivenessData, resp.StoreEnergyData, resp.TowerData, resp.SlashTowerData, resp.WeeklyData,
				resp.WeeklyRougeData,
			}
			renderItems = append(renderItems, resp.BattlePassData...)
		}

		signinTxt := "未签到"
		if resp.HasSignIn {
			signinTxt = "已签到"
		}
		title := []fyne.CanvasObject{
			widget.NewLabel(resp.ServerName + "--" + resp.RoleName),
			widget.NewLabel("签到: " + signinTxt),
		}

		return container.NewBorder(widget.NewSeparator(), widget.NewSeparator(),
			widget.NewSeparator(), widget.NewSeparator(),

			container.NewGridWithColumns(2,
				append(title, itemRender(renderItems...)...)...,
			),
		)
	}

	// 需要刷新组件
	changed := func() {
		tokenStr, err := vm.SelectedKujiequToken.Get()
		changedTime, _ := vm.RefreshKujiequWidget.Get()

		// 打开软件时不执行操作
		if changedTime == "" {
			return
		}
		if err != nil {
			slog.Info(tokenStr)
			return
		}
		kujiequWidgets := vm.GetKujiequWidgets(tokenStr)
		// kujiequWidgets := kujiequWidgetsFakeData
		// kujiequWidgets = append(kujiequWidgets, kujiequWidgetsFakeData...)
		kujiequWidgetsUI.RemoveAll()
		for _, wid := range kujiequWidgets {
			kujiequWidgetsUI.Add(card(wid))
		}
		kujiequWidgetsUI.Refresh()

	}

	// 添加监听事件
	vm.RefreshKujiequWidget.AddListener(binding.NewDataListener(changed))

	return kujiequLayout

}

// 添加token
func addTokenUI(win fyne.Window, vm *viewmodel.AppViewModel, onAdd func(token, devcode string)) fyne.CanvasObject {

	bindToken := binding.NewString()
	bindDevcode := binding.NewString()

	tokenEntry := widget.NewEntryWithData(bindToken)
	tokenEntry.SetPlaceHolder(T("输入token"))

	devcodeEntry := widget.NewEntryWithData(bindDevcode)
	devcodeEntry.SetPlaceHolder(T("devcode"))

	form := widget.NewForm(

		widget.NewFormItem(T("路径"), container.NewAdaptiveGrid(2, tokenEntry, devcodeEntry)),
	)
	form.SubmitText = T("添加")
	form.OnSubmit = func() {
		if err := form.Validate(); err != nil {
			slog.Info("form validate error", "err", err)
			return
		}
		token, _ := bindToken.Get()
		devcode, _ := bindDevcode.Get()
		onAdd(token, devcode)
		// 添加后重置表单状态
		bindToken.Set("")
		bindDevcode.Set("")
	}

	return widget.NewAccordion(
		widget.NewAccordionItem(T("添加Token"), form),
	)
}
