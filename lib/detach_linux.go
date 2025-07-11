package lib

import "syscall"

var (
	// linux 下进程分离
	exeDetachProcAttr = &syscall.SysProcAttr{
		Setsid: true, // For Linux, we use Setsid to detach the process
	}
)
