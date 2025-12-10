package genshin

// 原神切换服工具
import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type Lib struct {
	GameFile string `json:"game_file"`
	// 游戏路径
	GamePath string `json:"game_path"`

	// 存放配置文件与资源的路径
	ConfigPath string
	Log        func(...any)
}

func (lib *Lib) Init() {

}
func (lib *Lib) SetConfigPath(configPath string) {
	if configPath == "" {
		return
	}
	lib.ConfigPath = configPath
}

// 设置游戏可执行路径
func (lib *Lib) SetGameFile(gameFile string) {
	if gameFile == "" {
		return
	}
	lib.GameFile = gameFile
	lib.GamePath = filepath.Dir(gameFile)

}

// 更新游戏配置文件
func (lib *Lib) ServerConfig(serverName byte) {
	var modify [3]string
	switch serverName {
	case 'b':
		modify = [3]string{"14", "bilibili", "0"}
	case 'g':
		modify = [3]string{"1", "pcweb", "1"}

	default:
		lib.logInfo("暂不支持的server", serverName)
		return
	}
	// 将B服专用SDK复制到游戏目录下
	lib.cpBiliBiliSDK(serverName != 'b')

	changes := map[string]string{
		"channel":     modify[0],
		"cps":         modify[1],
		"sub_channel": modify[2],
		// b服的plugin version， 官服也能存在(但不使用)
		"plugin_sdk_version": "5.0.4",
	}

	gameConfigPath := filepath.Join(lib.GamePath, "config.ini")
	// 读取 INI 文件
	cfg, err := ini.Load(gameConfigPath)
	if err != nil {
		lib.logInfo("读取游戏配置文件失败：%v\n", err)
		return
	}

	// 获取或设置配置项的值
	section := cfg.Section("General")

	// 读取的配置和将要写入的配置是否相同
	hasDiff := false
	for k, v := range changes {
		if section.Key(k).String() != v {
			hasDiff = true
		}
		section.Key(k).SetValue(v)
	}

	// 配置文件与将要写入的内容没有不同，则直接返回
	if !hasDiff {
		lib.logInfo("游戏配置文件不需要更新")
		return
	}

	// 保存 INI 文件
	err = cfg.SaveTo(gameConfigPath)
	if err != nil {
		lib.logInfo("保存游戏配置文件失败：%v\n", err)
		return
	}

	lib.logInfo(gameConfigPath, "已更新并保存游戏配置文件")
}

// 目标为 B 服时如果游戏目录下没有SDK则拷贝SDK，为官服时若游戏目录下有 SDK 则删除
func (lib *Lib) cpBiliBiliSDK(remove bool) {
	// 将sdk移动到对应位置， 此sdk为b服专有
	sdkSourcePath := lib.ConfigPath
	sdkFileName := "PCGameSDK.dll"
	sourcePath := filepath.Join(sdkSourcePath, sdkFileName)
	targetPath := filepath.Join(lib.GamePath, "YuanShen_Data", "Plugins", sdkFileName)

	// 判断文件是否存在
	_, err := os.Stat(targetPath)
	sdkHasExist := os.IsExist(err)

	if remove && sdkHasExist {
		lib.logInfo("移除B服SDK文件")
		os.Remove(targetPath)
		return
	}

	// 需要copy SDK 且 SDK 不存在
	if !remove && !sdkHasExist {
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			// 读取错误处理
			lib.logInfo("读取 PCGameSDK.dll 文件失败, 请将此文件移动到")
			lib.logInfo("     " + sdkSourcePath)
		}
		err = os.WriteFile(targetPath, content, 0755)
		if err != nil {
			// 写文件出错处理
			lib.logInfo("SDK 写入游戏目录失败", err)
		}

		lib.logInfo("B服SDK已复制")
	}

}

// 将数据显示到界面中
func (lib *Lib) logInfo(args ...interface{}) {
	s := ""
	for _, arg := range args {
		s += fmt.Sprintf("%+v", arg)
	}
	if lib.Log != nil {
		lib.Log(s)
	}
	fmt.Println(s)

}

func NewLib() *Lib {
	lib := &Lib{}
	return lib
}
