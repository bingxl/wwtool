package viewmodel

import (
	"fmt"
	"log/slog"
	"slices"
	"wuwa/kujiequ"
	"wwtool/model"

	"fyne.io/fyne/v2/data/binding"
)

type KujiequViewModel struct {
	config                   *model.Config
	BindSelectedKujiequToken binding.String // 当前选中的库街区 token
	BindRefreshKujiequWidget binding.String // 刷新库街区widgets
}

func NewKujiequViewModel() (vm KujiequViewModel) {
	config, _ := model.LoadConfig()
	vm.config = config
	vm.BindRefreshKujiequWidget = binding.NewString()
	vm.BindSelectedKujiequToken = binding.NewString()
	if len(vm.config.Tokens) != 0 {
		vm.BindSelectedKujiequToken.Set(vm.config.Tokens[0].Token)
	}

	return
}

// 添加库街区token
func (vm *KujiequViewModel) AddToken(token, devcode string) {
	vm.config.Tokens = append(vm.config.Tokens, model.KujiequToken{Token: token, Devcode: devcode})
}

func (vm *KujiequViewModel) GetTokens() (tokens []string) {
	if vm.config.Tokens == nil {
		return
	}
	for _, t := range vm.config.Tokens {
		tokens = append(tokens, t.Token)
	}
	return
}

// 移除库街区token
func (vm *KujiequViewModel) DeleteToken(token string) {
	vm.config.Tokens = slices.DeleteFunc(vm.config.Tokens, func(t model.KujiequToken) bool {
		return t.Token == token
	})
}

// 从当前选中的库街区token 中实例化 kujiequ.KujieQu
func (vm *KujiequViewModel) getKujiequToken() (token kujiequ.Token, err error) {
	tk, err := vm.BindSelectedKujiequToken.Get()
	if tk == "" {
		err = fmt.Errorf("未选中token")
	}
	if err != nil {
		return
	}

	for _, t := range vm.config.Tokens {
		if t.Token == tk {
			token = kujiequ.Token(t)
			break
		}
	}

	return
}

// 获取库街区小组件
func (vm *KujiequViewModel) GetKujiequWidgets(token string, roleIds ...string) (widgets []kujiequ.WidgetResponseData) {

	tk, err := vm.getKujiequToken()
	if err != nil {
		slog.Error("实例化库街区失败" + err.Error())
		return
	}
	k := kujiequ.NewKujieQu(tk, nil)
	slog.Info("开始获取小组件")

	if len(roleIds) == 0 {
		widgets, _ = k.GetAllWidgets()
		return
	}
	for _, roleId := range roleIds {
		slog.Info("roleId" + roleId)
		// @TODO
	}

	return
}

// 库街区与角色签到
func (vm *KujiequViewModel) KujiequSignin() string {
	tk, err := vm.getKujiequToken()
	if err != nil {
		slog.Error("error to get kujiequ Token " + err.Error())
	}
	signinResult := kujiequ.StartSign([]kujiequ.Token{tk}, nil)
	return signinResult
}
