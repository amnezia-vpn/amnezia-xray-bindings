//go:build !windows

package main

import "syscall"

func shutdownFd(fd uintptr) {
	_ = syscall.Shutdown(int(fd), syscall.SHUT_RDWR)
}
