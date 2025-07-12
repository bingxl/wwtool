package lib

import (
	"log/slog"
	"os"
	"path/filepath"
)

// 检测src是否软链接到target
//
//	src - 源路径
//	target - 目标路径
//
// return:
//
//	bool - src 是否已软链接到target
//	error - 操作过程中出现的错误信息，若操作成功则返回nil。
func IsSymlinkTo(src, target string) (bool, error) {
	srcReal, err := filepath.EvalSymlinks(src)
	if err != nil {
		return false, err
	}
	path2Real, err := filepath.EvalSymlinks(target)
	if err != nil {
		return false, err
	}
	return srcReal == path2Real, nil
}

// CreateSymlink 将src软链接到target。如果src已经指向target，则不做任何操作。
//
// 如果target不存在，根据createTarget参数决定是否创建target目录。
//
// 参数:
//
//	src - 源路径，即软链接文件的路径。
//	target - 目标路径，软链接指向的路径。
//	createTarget - 当target不存在时，是否创建该目录。
//
// 返回值:
//
//	error - 操作过程中出现的错误信息，若操作成功则返回nil。
func CreateSymlink(src, target string, createTarget bool) error {
	// 检查src是否已经是指向target的软链接
	// 调用IsSymlinkTo函数判断，若已经指向则直接返回nil，不做后续操作
	if isSymlink, _ := IsSymlinkTo(src, target); isSymlink {
		return nil
	}

	// 若src不是指向target的软链接，先删除src，为创建新的软链接做准备
	os.Remove(src)

	// 检测target是否存在
	// 使用os.Stat函数获取文件信息，若返回的错误表示文件不存在，则进入后续处理
	if _, err := os.Stat(target); os.IsNotExist(err) {
		// 根据createTarget参数决定是否创建target目录
		if createTarget {
			// 创建target目录，使用os.MkdirAll递归创建目录，权限为0755
			err := os.MkdirAll(target, 0755)
			if err != nil {
				// 记录创建目录失败的错误信息
				slog.Error("Create target directory failed",
					"path2", target,
					"error", err,
				)
				return err
			}
		} else {
			// 若createTarget为false且target不存在，记录错误信息并返回错误
			slog.Error("Path does not exist",
				"path2", target,
				"error", err,
			)
			return err
		}
	}

	// 创建软链接，将src链接到target
	return os.Symlink(target, src)
}
