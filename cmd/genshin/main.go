// 编译为命令行模式

package main

import (
	"fmt"
	"wwtool/genshin"
)

var lib = genshin.NewLib()

func main() {
	lib.SetGameFile("D:\\program\\games\\miHoYo Launcher\\games\\Genshin Impact Game\\YuanShen.exe")
	fmt.Println("当前游戏路径: ", lib.GameFile)
	fmt.Println("输入操作后按回车\n g 切换官服 \n b 切换bilibili服 \n m path 修改游戏路径 eg: m D://xxxx/xx/Yuanshen.exe \n q  退出 \n s  启动游戏")

	var input string

	for {
		fmt.Scanln(&input)
		if input == "q" || input == "exit" {

			break
		}

		if len(input) > 2 && input[0] == 'm' && input[1] == ' ' {
			newPath := input[2:]
			if newPath == "" {
				fmt.Println("请输入有效路径")
				continue
			}
			lib.SetGameFile(newPath)
			fmt.Println("已修改游戏路径为: ", newPath)
			continue
		}
		operation := input[0]

		lib.ServerConfig(operation)

	}
}
