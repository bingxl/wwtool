package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TestUI() fyne.CanvasObject {

	// return InnerWindow()
	// 这里可以返回一个简单的测试 UI
	return container.NewVBox(

		// DocTab(),      // 添加文档标签页
		// InnerWindow(), // 添加内嵌窗口
		// Split(),
		ToolBar(), // 添加工具栏
	)
}

func DocTab() fyne.CanvasObject {
	// 这里可以返回一个简单的文档 UI
	tbs := []*container.TabItem{
		container.NewTabItem("Tab1", widget.NewLabel("这是第一个标签页")),
		container.NewTabItem("Tab2", widget.NewLabel("这是第二个标签页")),
	}
	doc := container.NewDocTabs(tbs...)
	// doc.SetTabLocation(container.TabLocationTrailing)

	return doc
}

func InnerWindow() fyne.CanvasObject {
	w1 := container.NewInnerWindow("inner window 1", widget.NewLabel("这是一个内嵌窗口"))
	w2 := container.NewInnerWindow("inner window 2", widget.NewLabel("这是另一个内嵌窗口"))
	return container.NewMultipleWindows(w1, w2)
}

func Split() fyne.CanvasObject {
	// 这里可以返回一个简单的分割视图
	left := widget.NewLabel("左侧内容")
	right := widget.NewLabel("右侧内容")
	split := container.NewHSplit(left, right)
	split.SetOffset(0.2) // 设置初始分割位置
	return split
}

func ToolBar() fyne.CanvasObject {
	label := widget.NewLabel("Toolbar Action Output")

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			label.SetText("Save clicked")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			label.SetText("Settings clicked")
		}),
	)

	content := container.NewVBox(
		toolbar,
		label,
	)
	return content
}
