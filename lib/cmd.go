package lib

import (
	"log/slog"
	"os/exec"
)

// 运行外部命令
// exe: 可执行文件路径
// detach: 是否分离进程
// 返回 *exec.Cmd 和 error
// 如果 detach 为 true，则在 Windows 上使用 DETACHED_PROCESS 标志分离进程
// 在 Linux 上使用 Setsid 分离进程
func RunExe(exe string, detach bool) (*exec.Cmd, error) {
	cmd := exec.Command(exe)
	if detach {
		cmd.SysProcAttr = exeDetachProcAttr
	}
	slog.Info("Running executable",
		"exe", exe,
		"detach", detach,
	)
	return cmd, cmd.Start()
}
