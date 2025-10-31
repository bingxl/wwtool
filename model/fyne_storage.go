package model

import (
	"fmt"
	"io"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

var (
	AppStorageRoot fyne.URI
	configFileName = "config.yaml"
	configFileURI  fyne.URI
)

func SetAppStorageRoot(root fyne.URI) {
	AppStorageRoot = root
	configFileURI, _ = storage.Child(root, configFileName)
}

func ConfigFileIsExists() bool {
	exists, err := storage.Exists(configFileURI)
	if err != nil {
		slog.Info("配置文件不存在", "configFileURI", configFileURI)
	}
	return exists
}

func ReadConfigFile() (data []byte, err error) {
	if !ConfigFileIsExists() {
		return nil, fmt.Errorf("配置文件不存在，或者不能读取")
	}
	canRead, err := storage.CanRead(configFileURI)
	if err != nil {
		return
	}
	if !canRead {
		return nil, fmt.Errorf("配置文件处于不可读状态")
	}

	reader, err := storage.Reader(configFileURI)
	if err != nil {
		return
	}
	defer reader.Close()
	data, err = io.ReadAll(reader)
	return
}

func WriteConfigToStorage(data []byte) (err error) {

	// Writer, 当资源不存在时会创建
	write, err := storage.Writer(configFileURI)
	if err != nil {
		return
	}
	defer write.Close()
	_, err = write.Write(data)
	return
}
