package viewmodel

import (
	"errors"
	"log/slog"
	"path/filepath"
	"wwtool/card"
	"wwtool/lib"
	"wwtool/model"

	"fyne.io/fyne/v2"
)

type AppViewModel struct {
	App          *fyne.App
	GamePathList []string
	SelectedID   int
	ConfigFile   string
	LinkServers  map[string]string // 软链接配置 {server:targetPath}

	// 配置项
	config *model.Config
}

func NewAppViewModel(configFile string, app *fyne.App) (*AppViewModel, error) {
	config, err := model.LoadConfig(configFile)
	if err != nil {
		return nil, err
	}

	return &AppViewModel{
		App:          app,
		GamePathList: config.GamePaths,
		SelectedID:   config.LastSelectedPath,
		config:       config,
		ConfigFile:   configFile,
		LinkServers:  config.LinkServers,
	}, nil
}

// 添加游戏路径
func (vm *AppViewModel) AddGamePath(path string) {
	vm.GamePathList = append(vm.GamePathList, path)
	vm.SaveConfig()
}

// 删除游戏路径
// index -1 表示删除当前选中项
// index >= 0 表示删除指定索引的项
// 如果索引超出范围，则不执行删除操作
func (vm *AppViewModel) RemoveGamePath(index int) {
	if index == -1 {
		index = vm.SelectedID
	}
	if index < 0 || index >= len(vm.GamePathList) {
		slog.Error("无效的索引", "index", index, "gamePathListLength", len(vm.GamePathList))
		return
	}
	vm.GamePathList = append(vm.GamePathList[:index], vm.GamePathList[index+1:]...)
	if vm.SelectedID >= len(vm.GamePathList) {
		vm.SelectedID = len(vm.GamePathList) - 1 // 如果删除后选中ID超出范围，调整为最后一个
	}

	vm.SaveConfig()
}

func (vm *AppViewModel) SetSelectedID(id int) {
	if id < 0 || id >= len(vm.GamePathList) {
		slog.Error("无效的选中ID", "id", id, "gamePathListLength", len(vm.GamePathList))
		return
	}
	vm.SelectedID = id
	vm.SaveConfig()
}

func (vm *AppViewModel) SetLinkServer(name string, path string) {
	if vm.LinkServers == nil {
		vm.LinkServers = make(map[string]string)
	}
	if path == "" {
		delete(vm.LinkServers, name)
	}
	vm.LinkServers[name] = path
	vm.SaveConfig()
}

// 保存配置文件
func (vm *AppViewModel) SaveConfig() error {

	vm.config.GamePaths = vm.GamePathList
	vm.config.LastSelectedPath = vm.SelectedID
	vm.config.LinkServers = vm.LinkServers

	return model.SaveConfig(vm.ConfigFile, vm.config)
}

func (vm *AppViewModel) RunSelected() error {
	if vm.SelectedID < 0 || vm.SelectedID >= len(vm.GamePathList) {
		slog.Error("无效的选中ID", "id", vm.SelectedID, "gamePathListLength", len(vm.GamePathList))
		return nil
	}
	exe := filepath.Join(vm.GamePathList[vm.SelectedID], vm.config.ExeName)
	_, err := lib.RunExe(exe, true)
	return err
}

func (vm *AppViewModel) CreateLinkToServer(server string) error {
	game := vm.GamePathList[vm.SelectedID]
	if game == "" {
		return errors.New("请先选择游戏路径")
	}
	targetPath, ok := vm.LinkServers[server]
	if !ok || targetPath == "" {
		return errors.New("未配置软链接目标路径")
	}
	src := filepath.Join(game, vm.config.LinkSrcPath)
	err := lib.CreateSymlink(src, targetPath)
	slog.Info("创建软链接", "server", server, "src", src, "targetPath", targetPath, "err", err)
	return err
}

// 获取游戏抽卡链接
func (vm *AppViewModel) GetGachaLink() string {
	link, err := card.GetLinkFromLog(vm.GamePathList[vm.SelectedID])
	vm.Clipboard(link)
	slog.Info("获取游戏抽卡链接", "link", link, "err", err)
	return link
}

func (vm *AppViewModel) Clipboard(src string) {
	(*vm.App).Clipboard().SetContent(src)
}
