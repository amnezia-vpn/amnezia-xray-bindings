//go:build windows

package main

import "syscall"

func shutdownFd(fd uintptr) {
	_ = syscall.Shutdown(syscall.Handle(fd), syscall.SHUT_RDWR)
}
