package model

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GamePaths        []string          `yaml:"game-paths"`         // 游戏路径列表, 存放游戏根目录路径，不是启动器路径
	ExeName          string            `yaml:"exe-name"`           // 游戏可执行文件名
	LinkServers      map[string]string `yaml:"link-servers"`       // 软链接配置 {server:targetPath}
	LastSelectedPath int               `yaml:"last-selected-path"` // 上次选中的路径索引
	LinkSrcPath      string            `yaml:"link-src-path"`      // 软链接源路径,相对于游戏跟目录
}

func LoadConfig(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return &Config{}, nil
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	return &config, err
}

func SaveConfig(file string, config *Config) error {
	data, err := yaml.Marshal(config)
	slog.Info("保存配置", "file", file, "config", config)
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}
