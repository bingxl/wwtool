package lib

import "syscall"

var (
	// Windows 下进程分离
	exeDetachProcAttr = &syscall.SysProcAttr{
		// CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | DETACHED_PROCESS,
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
)
