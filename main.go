package main

import (
	"log"
	"log/slog"

	"os"

	"wwtool/i18n"
	"wwtool/model"
	"wwtool/view"
	"wwtool/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

const (
	appID = "wwtool"
)

func main() {
	// 访问 查看内存分析 http://localhost:6060/debug/pprof/heap
	// go func() {
	// 	http.ListenAndServe("localhost:6060", nil)
	// }()

	// 设置多语言
	// langs := i18n.GetSupportLang()
	// slog.Info("support langs", "langs", langs)
	// i18n.SetLang(langs[1])
	// 初始化配置文件（首次创建）
	model.SetUniqueID(appID)

	// 使用指定的字体
	// setFont()
	a := app.NewWithID(appID)
	w := a.NewWindow("ww tool")

	// 设置应用的生命周期事件处理函数
	viewmodel.SetLifecycle(&a)
	// setMainMenu(w)

	vm, err := viewmodel.NewAppViewModel()
	if err != nil {
		log.Fatal(err)
	}
	w.Resize(fyne.NewSize(690, 500))
	// w.SetFixedSize(true)
	w.SetContent(view.BuildUI(w, vm))

	w.ShowAndRun()
}

func setFont() {
	os.Setenv("FYNE_FONT", "C:/Windows/Fonts/Deng.ttf")
}

func setMainMenu(w fyne.Window) {
	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu(i18n.T("文件"),
			fyne.NewMenuItem(i18n.T("退出"), func() {
				slog.Info("退出应用")
			}),
		),
		fyne.NewMenu(i18n.T("编辑"),
			fyne.NewMenuItem(i18n.T("首选项"), func() {
				slog.Info("打开首选项")
			}),
		),
	)
	w.SetMainMenu(mainMenu)
}
