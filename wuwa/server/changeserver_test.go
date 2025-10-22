package server

import (
	"log/slog"
	"testing"
)

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

func TestChangeServer(t *testing.T) {
	// 这里是测试代码
	// 你可以使用 ChangeServer 的方法进行测试
	// 例如：
	cs := ChangeServer{sourcePath: "../tmp/source", BiliPath: "../tmp/bilibili", OfficialPath: "../tmp/official"}

	err := cs.createLink(cs.OfficialPath)
	if err != nil {
		t.Error("Error:", err)
	} else {
		t.Log("Changed to Bilibili successfully")
	}
}
