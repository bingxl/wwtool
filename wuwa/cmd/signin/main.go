package main

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"wuwa"
	"wuwa/kujiequ"
)

var (
	configFile = flag.String("config", "", "配置文件")
	command    = flag.String("cmd", "signin", "运行的指令， signin | widget")
)

// powershell  $env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags="-s -w" -o tmp/kurobbsautosignin .\cmd\main.go
// go run .\cmd\signin\main.go -cmd=widget -config="./config.json"
func main() {
	// startCardRecodeRequest()
	flag.Parse()

	if *configFile == "" {
		dir, err := os.Executable()
		if err != nil {
			slog.Error("获取可执行文件目录失败" + err.Error())
			dir = "./"
		}

		dir = filepath.Dir(dir)
		*configFile = filepath.Join(dir, "config.json")
	}

	config, err := wuwa.GetConfig(*configFile)
	if err != nil {
		slog.Info("获取config失败" + err.Error())
		return
	}
	if len(config.Tokens) <= 0 {
		slog.Error("不存在token")
		return
	}

	switch *command {
	case "signin":
		str := kujiequ.StartSign(config.Tokens, config.KujiequHeaders)
		slog.Info("签到结果", "result", str)
	case "widget":
		getWidget(config)
	}

	// card.StartCardRecodeRequest(config.WuwaGamePath[0])
}

func getWidget(config wuwa.Config) {
	for _, token := range config.Tokens {
		k := kujiequ.NewKujieQu(token, config.KujiequHeaders)

		roles, err := k.FindAllRoles()
		if err != nil {
			slog.Error("获取角色信息失败", "token", token, "err", err)
			continue
		}
		for _, role := range roles {
			widget, err := k.GetWidget(role)
			if err != nil {
				slog.Error("获取小组件失败", "角色名", role.RoleName, "err", err)
				continue
			}
			slog.Info("", "widget", widget)
		}

	}
}
