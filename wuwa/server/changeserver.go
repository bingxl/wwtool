package server

import (
	"errors"
	"log/slog"
	"os"
)

// b服官服 切换

var OfficialName = "official"
var BiliName = "bilibili"

var ()

type ChangeServer struct {
	sourcePath   string
	BiliPath     string
	OfficialPath string
}

func (c ChangeServer) Check() error {
	return errors.New("not implemented")
}

// 根据传入的serverName 创建软链接，使sourcePath指向BiliPath或OfficialPath
func (c ChangeServer) Change(serverName string) error {
	var targetPath = ""
	switch serverName {
	case BiliName:
		targetPath = c.BiliPath
	case OfficialName:
		targetPath = c.OfficialPath
	default:
		return errors.New("not implement: " + serverName)
	}

	return c.createLink(targetPath)

}

func (c ChangeServer) createLink(targetPath string) error {
	slog.Debug("接收到的TargetPath:" + targetPath)
	// 检查目标路径是否存在
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		errs := errors.Join(err, errors.New("软链接目标路径不存在:"+targetPath))
		slog.Error(errs.Error())
		return errs
	}

	// 删除原有的软链接
	if err := os.RemoveAll(c.sourcePath); err != nil {
		slog.Error(err.Error())
		return err
	}
	// 创建新的软链接
	return os.Symlink(targetPath, c.sourcePath)

}
