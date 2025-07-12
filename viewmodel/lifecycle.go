package viewmodel

import (
	"fyne.io/fyne/v2"
)

// fyne 应用生命周期事件处理函数

type EventName string

const (
	EnteredForeground EventName = "onEnteredForeground"
	ExitedForeground  EventName = "onExitedForeground"
	Started           EventName = "onStarted"
	Stopped           EventName = "onStopped"
)

var (
	events = map[EventName]([]func()){}
)

func runEvent(event EventName) {
	for _, fn := range events[event] {
		fn()
	}
}

// SetLifecycle 设置应用的生命周期事件处理函数
func SetLifecycle(app *fyne.App) {
	cycle := (*app).Lifecycle()
	cycle.SetOnEnteredForeground(func() { runEvent(EnteredForeground) })
	cycle.SetOnExitedForeground(func() { runEvent(ExitedForeground) })
	cycle.SetOnStarted(func() { runEvent(Started) })
	cycle.SetOnStopped(func() { runEvent(Stopped) })
}

// RegisterEvent 注册一个生命周期事件处理函数
func RegisterEvent(event EventName, fn func()) bool {
	if _, ok := events[event]; !ok {
		// slog.Error("未知的生命周期事件", "event", event)
		events[event] = []func(){}
	}
	events[event] = append(events[event], fn)
	return true
}

// UnregisterEvent 注销一个生命周期事件处理函数
func RemoveEvent(event EventName, fn func()) {
	if _, ok := events[event]; ok {
		for i, f := range events[event] {
			if &f == &fn {
				events[event] = append(events[event][:i], events[event][i+1:]...)
				break
			}
		}
	}
}
