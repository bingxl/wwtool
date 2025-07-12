package main

import (
	"log"
	"log/slog"

	"os"

	"wwtool/model"
	"wwtool/view"
	"wwtool/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

const (
	appID = "wwtool"
)

func main() {
	// 访问 查看内存分析 http://localhost:6060/debug/pprof/heap
	// go func() {
	// 	http.ListenAndServe("localhost:6060", nil)
	// }()

	// 初始化配置文件（首次创建）
	model.SetUniqueID(appID)

	// 使用指定的字体
	// setFont()
	a := app.NewWithID(appID)
	w := a.NewWindow("ww tool")

	// 设置应用的生命周期事件处理函数
	viewmodel.SetLifecycle(&a)
	// setMainMenu(w)

	vm, err := viewmodel.NewAppViewModel(&a)
	if err != nil {
		log.Fatal(err)
	}

	w.SetContent(container.NewVBox(view.BuildUI(w, vm)))
	w.Resize(fyne.NewSize(500, 400))
	w.ShowAndRun()
}

func setFont() {
	os.Setenv("FYNE_FONT", "C:/Windows/Fonts/Deng.ttf")
}

func setMainMenu(w fyne.Window) {
	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Quit", func() {
				slog.Info("退出应用")
			}),
		),
		fyne.NewMenu("Edit",
			fyne.NewMenuItem("Preferences", func() {
				slog.Info("打开首选项")
			}),
		),
	)
	w.SetMainMenu(mainMenu)
}
