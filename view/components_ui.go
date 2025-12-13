package view

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Title
//
// 传入的Title 会自动调用 T 翻译
func Title(title string) fyne.CanvasObject {
	return widget.NewLabelWithStyle(
		T(title),
		fyne.TextAlignLeading,
		fyne.TextStyle{
			Bold:   true,
			Italic: true,
		},
	)
}

// ShowFolderOpen 弹出文件夹选择对话框，返回用户选择的文件夹路径或错误
// 如果用户取消选择，则返回空字符串和 nil 错误
func ShowFolderOpen(win fyne.Window) (folder string, errReturn error) {
	ch := make(chan bool)
	go func() {
		dialog.ShowFolderOpen(func(reader fyne.ListableURI, err error) {
			defer close(ch) // 确保通道在函数结束时被关闭

			if err != nil {
				errReturn = err
				return
			}
			if reader == nil {
				folder = ""
				errReturn = nil // 用户取消选择文件夹
				return          // 用户取消选择文件夹
			}
			folder = reader.Path()
			if folder == "" {
				return // 无效路径
			}

		}, win)
	}()

	<-ch // 等待用户选择完成
	return
}

// 定时关闭的info
func ShowInfoWithAutoClose(title, message string, win fyne.Window) {
	dia := dialog.NewInformation(title, message, win)
	dia.Show()
	time.AfterFunc(2*time.Second, func() {
		fyne.Do(dia.Hide)
	})
}
