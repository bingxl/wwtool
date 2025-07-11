package lib

import (
	"log/slog"
	"os"
	"path/filepath"
)

// 检测src是否软链接到path2
func IsSymlinkTo(src, path2 string) (bool, error) {
	// 获取path1的软链接目标
	target, err := os.Readlink(src)
	if err != nil {
		slog.Error("Failed to read symlink",
			"src", src,
			"error", err,
		)
		return false, err
	}

	// 获取绝对路径以便比较
	absTarget, err := filepath.Abs(target)
	if err != nil {
		slog.Error("Failed to get absolute path of symlink target",
			"target", target,
			"error", err,
		)
		return false, err
	}
	absPath2, err := filepath.Abs(path2)
	if err != nil {
		slog.Error("Failed to get absolute path of path2",
			"path2", path2,
			"error", err,
		)
		return false, err
	}
	// 比较目标路径是否匹配path2或path3
	return absTarget == absPath2, nil
}

// 将src软链接到path2
func CreateSymlink(src, path2 string) error {
	// 检查src是否已存在
	if _, err := os.Lstat(src); err == nil {
		os.Remove(src) // 如果存在则删除
	}

	// 检测path2是否存在，不存在则返回错误
	if _, err := os.Stat(path2); os.IsNotExist(err) {
		slog.Error("Path does not exist",
			"path2", path2,
			"error", err,
		)
		return err
	}

	// 创建软链接
	return os.Symlink(path2, src)
}
