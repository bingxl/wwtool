package viewmodel

import (
	"path/filepath"
	"slices"
	"wwtool/lib"
	"wwtool/model"

	"fyne.io/fyne/v2/data/binding"
)

type ThirdToolsViewModel struct {
	config         *model.Config
	BindThirdNames binding.StringList
}

// ThirdTools 相关操作
//
// 包含增加/删除/获取/运行
func (vm *ThirdToolsViewModel) GetThirdToolsName() []string {
	names := make([]string, len(vm.config.ThirdTools))
	for i, v := range slices.All(vm.config.ThirdTools) {
		names[i] = filepath.Base(v)
	}
	return names
}
func (vm *ThirdToolsViewModel) FreshBindNames() {
	names := vm.GetThirdToolsName()
	vm.BindThirdNames.Set(names)
}

func (vm *ThirdToolsViewModel) AddThirdTools(pa string) {
	if slices.Contains(vm.config.ThirdTools, pa) {
		return
	}
	vm.config.ThirdTools = append(vm.config.ThirdTools, pa)
	vm.FreshBindNames()
}
func (vm *ThirdToolsViewModel) RemoveThirdTools(index int) {
	if index >= 0 && index < len(vm.config.ThirdTools) {
		vm.config.ThirdTools = slices.Delete(vm.config.ThirdTools, index, index+1)
		vm.FreshBindNames()
	}

}
func (vm *ThirdToolsViewModel) RunExe(index int) {
	if 0 <= index && index < len(vm.config.ThirdTools) {
		lib.RunExe(vm.config.ThirdTools[index], true)
	}

}

func NewThirdToolsViewModel() (tool ThirdToolsViewModel) {
	// config 是全局单列模式
	config, _ := model.LoadConfig()
	tool.config = config
	tool.BindThirdNames = binding.NewStringList()
	tool.BindThirdNames.Set(tool.GetThirdToolsName())

	return
}
