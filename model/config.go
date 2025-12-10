package model

import (
	"embed"
	"log/slog"

	"gopkg.in/yaml.v3"
)

// 默认配置文件
//
//go:embed default_config.yaml
var defaultConfig embed.FS

type KujiequToken struct {
	Token   string `yaml:"token"`
	Devcode string `yaml:"devcode"`
}

type GamePath struct {
	Name              string `yaml:"name"`                // 自定义的名称
	Path              string `yaml:"path"`                // 游戏所在目录
	HasMultipleServer bool   `yaml:"has-multiple-server"` // 游戏本体是否能用于大陆官服/b服

}

type Config struct {
	GamePaths        []GamePath        `yaml:"game-paths"`         // 游戏路径列表, 存放游戏根目录路径，不是启动器路径
	ExeName          string            `yaml:"exe-name"`           // 游戏可执行文件名
	LinkServers      map[string]string `yaml:"link-servers"`       // 软链接配置 {server:targetPath}
	LastSelectedPath int               `yaml:"last-selected-path"` // 上次选中的路径索引
	LinkSrcPath      string            `yaml:"link-src-path"`      // 软链接源路径,相对于游戏跟目录
	ThirdTools       []string          `yaml:"third-tools"`        // 小工具
	Tokens           []KujiequToken    `yaml:"tokens"`             // 库街区token

	GenshinPath string `yaml:"genshin-path"` // 原神游戏路径
}

var (
	// 全局配置
	globalConfig *Config
	appUniqueID  string = "wwtool"

	// 配置文件存放路径
	configFile string
)

func SetUniqueID(uid string) {
	appUniqueID = uid
}

// 加载默认配置
func loadDefault() *Config {
	slog.Info("load default config")
	cfg := &Config{}
	data, err := defaultConfig.ReadFile("default_config.yaml")
	if err == nil {
		yaml.Unmarshal(data, &cfg)
	}
	slog.Debug("in loadDefault function", "defaultConfig", cfg)
	return cfg
}

// 获取Config
//
// 已加载过则直接返回， 否则读取配置文件返回， 如果配置文件不存在则返回默认配置

func LoadConfig() (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}
	slog.Debug("config fyne.URI ", "name", configFileURI.Name(), "string", configFileURI.String(),
		"path", configFileURI.Path(), "extension", configFileURI.Extension(),
		"scheme", configFileURI.Scheme(),
	)
	// 配置文件不存在，使用默认配置
	if !ConfigFileIsExists() {
		globalConfig = loadDefault()
		return globalConfig, nil
	}

	// 读取配置文件
	data, err := ReadConfigFile()
	if err != nil {
		return &Config{}, nil
	}
	err = yaml.Unmarshal(data, &globalConfig)
	if err != nil {
		slog.Error("解析配置文件失败, 使用默认配置", "err", err)
		globalConfig = loadDefault()
	}
	return globalConfig, nil
}

// 保存全局配置
func SaveConfig() error {
	data, err := yaml.Marshal(globalConfig)
	slog.Info("保存配置", "file", configFile, "config", globalConfig)
	if err != nil {
		return err
	}

	return WriteConfigToStorage(data)
}
