package viewmodel

import (
	"wwtool/genshin"
	"wwtool/model"

	"fyne.io/fyne/v2/data/binding"
)

type GenshinViewModel struct {
	config      *model.Config
	genshinLib  genshin.Lib
	GenshinPath binding.String
}

// 设置原神游戏路径
func (vm *GenshinViewModel) SetGenshinPath(name string) {
	vm.GenshinPath.Set(name)
	vm.config.GenshinPath = name
	vm.genshinLib.SetGameFile(name)
}

// 切换原神服务器
func (vm *GenshinViewModel) SwitchServer(server byte) {
	vm.genshinLib.ServerConfig(server)
}

func NewGenshinViewModel() *GenshinViewModel {
	config, _ := model.LoadConfig()

	gVM := GenshinViewModel{
		config:      config,
		GenshinPath: binding.NewString(),
		genshinLib:  *genshin.NewLib(),
	}

	if len(config.GenshinPath) > 0 {
		gVM.GenshinPath.Set(config.GenshinPath)
	}
	gVM.genshinLib.SetGameFile(config.GenshinPath)

	return &gVM
}
