package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"wwtool/card"
	"wwtool/lib"
)

var (
	// 目标路径映射
	serverPaths = map[string]string{
		"cn":   `C:\program\games\wuwa-diff\KrPcSdk_Mainland-official`,
		"bili": `C:\program\games\wuwa-diff\KrPcSdk_Mainland-bili`,
	}
	// 链接路径
	linkPath = `C:\program\games\Wuthering Waves\Wuthering Waves Game\Client\Binaries\Win64\ThirdParty\KrPcSdk_Mainland`
	gamePath = `C:\program\games\Wuthering Waves\Wuthering Waves Game` // 游戏本体目录
)

var operationMap = map[string]func(){
	// 切换到官服
	"1": func() { changeLink("cn") },
	// 切换到b服
	"2": func() { changeLink("bili") },

	// 输出抽卡链接
	"3": func() {
		link, err := card.GetLinkFromLog(gamePath)
		slog.Info("抽卡链接", "link", link, "error", err)
	},

	// 退出程序
	"q": func() { os.Exit(0) },
	"t": test,
}

func main() {
	curPath, _ := os.Executable()
	slog.Info("cur path", "cur", curPath)
	server := "cn"
	if len(os.Args) > 1 {
		server = os.Args[1]
		changeLink(server)
	}
	for {
		operation()
	}
}

// 创建软链接
func changeLink(server string) {
	targetPath, exists := serverPaths[server]
	if !exists {
		slog.Error("Invalid server selection", "server", server)
		return
	}

	slog.Info("Changing symlink",
		"linkPath", linkPath,
		"targetPath", targetPath,
	)

	// 创建新的软链接
	err := lib.CreateSymlink(linkPath, targetPath, false)
	if err != nil {
		slog.Error("Failed to create symlink", "error", err)
	} else {
		slog.Info("change server to " + server)
	}
}

// 获取目标服务器名称
func operation() {
	// 交互式选择
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("请选择要执行的操作：")
	fmt.Println("1. 切换到官服")
	fmt.Println("2. 切换到B服")
	fmt.Println("3. 输出抽卡链接")
	fmt.Println("t. run test function")
	fmt.Println("q. 退出程序")
	fmt.Print("输入数字选择 (默认1): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if _, ok := operationMap[input]; ok {
		operationMap[input]()
	} else if input == "" {
		// 如果没有输入，默认选择1
		operationMap["1"]()
	}
}

func test() {
	src := "path1"
	path2 := "path2"
	path3 := "path3"

	slog.Info("Checking if symlink exists",
		"src", src,
		"path2", path2,
		"path3", path3,
	)
	exists, err := lib.IsSymlinkTo(src, path2)
	slog.Info("Symlink check result",
		"src", src,
		"path2", path2,
		"exists", exists,
		"error", err,
	)
	fmt.Println("creating symlink")
	fmt.Println(lib.CreateSymlink(src, path2, false))
	real, err := filepath.EvalSymlinks(src)
	slog.Info("real path",
		"src", src,
		"path2", path2,
		"real", real,
	)
	real, err = filepath.EvalSymlinks("pathtest")
	slog.Info("real path",
		"path2", path2,
		"real", real,
		"error", err,
	)

	// fmt.Println("creating symlink from", src, "to", path3)
	// fmt.Println(lib.CreateSymlink(src, path3))

	lib.CreateSymlink(src, "path4", false)

	preDefinePath := map[string]func() (string, error){
		"userConfigDir": os.UserConfigDir,
		"userCacheDir":  os.UserCacheDir,
		"userHomeDir":   os.UserHomeDir,
		"TempDir":       func() (string, error) { return os.TempDir(), nil },
	}
	for key, v := range preDefinePath {
		dir, _ := v()
		slog.Info("预定义路径", key, dir)
	}
}
