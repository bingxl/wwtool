package lib

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	dir    = filepath.Join(os.TempDir(), "pathlink_test")
	src    = filepath.Join(dir, "src")
	target = filepath.Join(dir, "path2")
)

func isExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func rm(path string) error {
	if !isExist(path) {
		return nil
	}
	return os.Remove(path)
}
func mk(path string) error {
	if !isExist(path) {
		return os.Mkdir(path, 0755)
	}
	return nil
}

func TestCreateSymlink(t *testing.T) {
	type tCase struct {
		pre  func()
		info string
		run  func(t *testing.T)
	}
	testCase := []tCase{
		{
			pre:  func() { rm(target) },
			info: "target 不存在时",
			run: func(t *testing.T) {
				err := CreateSymlink(src, target, false)
				if err == nil {
					t.Errorf("createTarget == false 时 应该创建失败")
				} else {
					t.Log("createTarget == false 时测试成功")
				}

				rm(target)
				err = CreateSymlink(src, target, true)
				if err != nil {
					t.Errorf("createTarget == true 时应该创建成功")
				} else {
					t.Log("createTarget == true 时成功")
				}

			},
		},
		{
			pre:  func() { mk(target) },
			info: "target 存在时",
			run: func(t *testing.T) {
				err := CreateSymlink(src, target, false)
				if err != nil {
					t.Errorf("createTarget == false 时 应该创建成功")
				} else {
					t.Log("createTarget == false 时测试成功")
				}

				mk(target)
				err = CreateSymlink(src, target, true)
				if err != nil {
					t.Errorf("createTarget == true 时应该创建成功")
				} else {
					t.Log("createTarget == true 时成功")
				}

			},
		},
	}

	for _, tt := range testCase {
		tt.pre()
		t.Run(tt.info, tt.run)
	}

	// 清理测试产生的文件
	rm(target)
	rm(src)
	rm(dir)
}
