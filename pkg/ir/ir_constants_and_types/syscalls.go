// Package ir_constants_and_types содержит определения номеров системных вызовов,
// архитектурных регистров и типов данных для x86_64.
package ir_constants_and_types

// Константы номеров системных вызовов для архитектуры Linux x86_64.
const (
	SYSCALL_OPEN      = 2
	SYSCALL_CLOSE     = 3
	SYSCALL_SEND_FILE = 40
	SYSCALL_LSEEK     = 8
	SYSCALL_WRITE     = 1
	SYSCALL_EXIT      = 60
	SYSCALL_FORK      = 57
	SYSCALL_EXECVE    = 59
	SYSCALL_CHMOD     = 90
	SYSCALL_RENAME    = 82
	SYSCALL_UNLINK    = 87
	SYSCALL_NANOSLEEP = 35
	SYSCALL_GETPID    = 39
	SYSCALL_GETUID    = 102
	SYSCALL_GETGID    = 104
	SYSCALL_GETPPID   = 110
)

// SyscallNames сопоставляет числовые идентификаторы системных вызовов с их строковыми именами.
var SyscallNames = map[int64]string{
	SYSCALL_OPEN:      "sys_open",
	SYSCALL_CLOSE:     "sys_close",
	SYSCALL_SEND_FILE: "sys_sendfile",
	SYSCALL_LSEEK:     "sys_lseek",
	SYSCALL_WRITE:     "sys_write",
	SYSCALL_EXIT:      "sys_exit",
	SYSCALL_FORK:      "sys_fork",
	SYSCALL_EXECVE:    "sys_execve",
	SYSCALL_CHMOD:     "sys_chmod",
	SYSCALL_RENAME:    "sys_rename",
	SYSCALL_UNLINK:    "sys_unlink",
	SYSCALL_NANOSLEEP: "sys_nanosleep",
	SYSCALL_GETPID:    "sys_getpid",
	SYSCALL_GETUID:    "sys_getuid",
	SYSCALL_GETGID:    "sys_getgid",
	SYSCALL_GETPPID:   "sys_getppid",
}
