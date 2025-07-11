package main

import (
	"embed"
	"log"
	"log/slog"

	"os"

	"wwtool/view"
	"wwtool/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

//go:embed default_config.yaml
var defaultConfig embed.FS

const (
	appID = "wwtool"
)

func main() {
	// 访问 查看内存分析 http://localhost:6060/debug/pprof/heap
	// go func() {
	// 	http.ListenAndServe("localhost:6060", nil)
	// }()

	// 初始化配置文件（首次创建）
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		data, _ := defaultConfig.ReadFile("default_config.yaml")
		os.WriteFile("config.yaml", data, 0644)
	}
	setFont()
	a := app.NewWithID(appID)
	w := a.NewWindow("ww tool")

	// setLifecycle(a)
	// setMainMenu(w)

	vm, err := viewmodel.NewAppViewModel("config.yaml", &a)
	if err != nil {
		log.Fatal(err)
	}

	w.SetContent(container.NewVBox(view.BuildUI(w, vm)))
	w.Resize(fyne.NewSize(500, 400))
	w.ShowAndRun()
}

func setFont() {
	// os.Setenv("FYNE_FONT", "C:/Windows/Fonts/Deng.ttf")
}

func setLifecycle(a fyne.App) {
	cycle := a.Lifecycle()
	cycle.SetOnEnteredForeground(func() {
		slog.Info("应用已进入前台")
	})
	cycle.SetOnExitedForeground(func() {
		slog.Info("应用已进入后台")
	})
	cycle.SetOnStarted(func() {
		slog.Info("应用已启动")
	})
	cycle.SetOnStopped(func() {
		slog.Info("应用已停止")
	})
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
