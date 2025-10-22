package main

import (
	"log/slog"
	"wuwa/card"
)

func main() {
	// gamePath := "C:\\program\\games\\wuwa-overseas\\Wuthering Waves Game"
	gamePath := "C:\\program\\games\\Wuthering Waves\\Wuthering Waves Game"
	link, err := card.GetLinkFromLog(gamePath)
	slog.Info("result", "link", link, "err", err)
}
