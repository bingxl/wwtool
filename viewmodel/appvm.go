package viewmodel

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"wwtool/card"
	"wwtool/i18n"
	"wwtool/lib"
	"wwtool/model"

	"fyne.io/fyne/v2"
)

var T = i18n.T

type AppViewModel struct {
	App          *fyne.App
	GamePathList []model.GamePath

	SelectedID  int               // 当前选中的gamePath index， -1 表示都没选中
	LinkServers map[string]string // 软链接配置 {server:targetPath}

	// 配置项
	config *model.Config
}

func NewAppViewModel(app *fyne.App) (*AppViewModel, error) {
	config, err := model.LoadConfig()
	if err != nil {
		slog.Error("load config error", "config", config)
		return nil, err
	}
	slog.Info("config content", "config", config)
	vm := &AppViewModel{
		App:          app,
		GamePathList: config.GamePaths,
		SelectedID:   config.LastSelectedPath,
		config:       config,

		LinkServers: config.LinkServers,
	}
	RegisterEvent(Stopped, func() { vm.SaveConfig() })

	return vm, nil
}

// 添加游戏路径
//
// 如果name已存在则替换
func (vm *AppViewModel) AddGamePath(name, path string, hasMultipleServer bool) {
	if name == "" || path == "" {
		slog.Error("添加game 时 name与path都不能为空", "name", name, "path", path)
		return
	}
	game := model.GamePath{
		Name:              name,
		Path:              path,
		HasMultipleServer: hasMultipleServer,
	}
	if index := vm.GetIndexByGameName(name); index != -1 {
		vm.GamePathList[index] = game
	} else {
		vm.GamePathList = append(vm.GamePathList, game)
	}
}

// 删除游戏路径
//
//	name - 游戏路径的自定义名称
//
// 如果索引超出范围，则不执行删除操作
func (vm *AppViewModel) RemoveGamePath(name string) {
	vm.GamePathList = slices.DeleteFunc(
		vm.GamePathList,
		func(e model.GamePath) bool {
			return e.Name == name
		},
	)
}

// 返回由 gamePath.Name 组成的slice
func (vm *AppViewModel) GetGameNames() []string {
	names := make([]string, len(vm.GamePathList))
	for index, game := range vm.GamePathList {
		names[index] = game.Name
	}
	return names
}

// 通过gamePath.name 获取索引
func (vm *AppViewModel) GetIndexByGameName(name string) int {
	index := slices.IndexFunc(vm.GamePathList, func(game model.GamePath) bool {
		return game.Name == name
	})
	return index
}

// 设置当前选中的gamePath 索引
func (vm *AppViewModel) SetSelectedID(name string) {
	vm.SelectedID = vm.GetIndexByGameName(name)
}

func (vm *AppViewModel) SetLinkServer(name string, path string) {
	if vm.LinkServers == nil {
		vm.LinkServers = make(map[string]string)
	}
	if path == "" {
		delete(vm.LinkServers, name)
	}
	vm.LinkServers[name] = path
}

// 保存配置文件
func (vm *AppViewModel) SaveConfig() error {
	vm.config.GamePaths = vm.GamePathList
	vm.config.LastSelectedPath = vm.SelectedID
	vm.config.LinkServers = vm.LinkServers
	slog.Info("保存配置文件")

	return model.SaveConfig()
}

// 启动选中的程序
func (vm *AppViewModel) RunSelected() error {
	if !vm.ValidSelectedID() {
		return nil
	}
	path := vm.GamePathList[vm.SelectedID].Path
	exe := filepath.Join(path, vm.config.ExeName)
	_, err := lib.RunExe(exe, true)
	return err
}

// 判断vm.SelectedID 是否有效
func (vm *AppViewModel) ValidSelectedID() bool {
	if vm.SelectedID < 0 || vm.SelectedID >= len(vm.GamePathList) {
		slog.Error("无效的选中ID", "id", vm.SelectedID, "gamePathListLength", len(vm.GamePathList))
		return false
	}
	return true
}

// 创建软链接到软件缓存路径下
//
// 将游戏的src 目录软链接到 userCacheDir/uniqueID/serverPath 目录
func (vm *AppViewModel) CreateLinkToServer(server string) error {
	if !vm.ValidSelectedID() {
		return errors.New(T("游戏索引无效"))
	}
	game := vm.GamePathList[vm.SelectedID]
	if !game.HasMultipleServer {
		return errors.New(T("当前游戏不支持多个服"))
	}

	targetPath, ok := vm.LinkServers[server]
	if !ok || targetPath == "" {
		return errors.New(T("未配置软链接目标目录"))
	}
	src := filepath.Join(game.Path, vm.config.LinkSrcPath)

	// target 不存在时创建
	err := lib.CreateSymlink(src, filepath.Join(vm.GetCacheDirWithAppID(), targetPath), true)
	slog.Info("创建软链接", "server", server, "src", src, "targetPath", targetPath, "err", err)
	return err
}

// 获取软件缓存目录
//
// Windows 下为%AppData%/Local/%UniqueID%
func (vm *AppViewModel) GetCacheDirWithAppID() string {
	dir, _ := os.UserCacheDir()
	return filepath.Join(dir, (*vm.App).UniqueID())
}

// 获取游戏抽卡链接
func (vm *AppViewModel) GetGachaLink() string {
	if !vm.ValidSelectedID() {
		return ""
	}
	link, err := card.GetLinkFromLog(vm.GamePathList[vm.SelectedID].Path)
	vm.Clipboard(link)
	slog.Info("获取游戏抽卡链接", "link", link, "err", err)
	return link
}

// 复制content 到剪切板
func (vm *AppViewModel) Clipboard(content string) {
	(*vm.App).Clipboard().SetContent(content)
}

// ThirdTools 相关操作
//
// 包含增加/删除/获取/运行
func (vm *AppViewModel) GetThirdToolsName() []string {
	names := make([]string, len(vm.config.ThirdTools))
	for i, v := range slices.All(vm.config.ThirdTools) {
		names[i] = filepath.Base(v)
	}
	return names
}
func (vm *AppViewModel) AddThirdTools(pa string) {
	if slices.Contains(vm.config.ThirdTools, pa) {
		return
	}
	vm.config.ThirdTools = append(vm.config.ThirdTools, pa)
}
func (vm *AppViewModel) RemoveThirdTools(pa string) {
	vm.config.ThirdTools = slices.DeleteFunc(vm.config.ThirdTools, func(v string) bool {
		return pa == v
	})
}
func (vm *AppViewModel) RunExe(index int) {
	if 0 <= index && index < len(vm.config.ThirdTools) {
		lib.RunExe(vm.config.ThirdTools[index], true)
	}

}
