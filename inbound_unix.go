//go:build !windows

package main

import "syscall"

// Accepted inbound sockets are not reachable through xray-core, so after the
// listeners are closed every connected socket bound to an inbound port is
// shut down by walking the process file descriptors.
func shutdownInboundConns(ports map[int]struct{}) {
	if len(ports) == 0 {
		return
	}

	maxFd := 65536
	var rl syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rl); err == nil && rl.Cur < uint64(maxFd) {
		maxFd = int(rl.Cur)
	}

	for fd := 0; fd < maxFd; fd++ {
		local, err := syscall.Getsockname(fd)
		if err != nil {
			continue
		}
		var port int
		switch a := local.(type) {
		case *syscall.SockaddrInet4:
			port = a.Port
		case *syscall.SockaddrInet6:
			port = a.Port
		default:
			continue
		}
		if _, ok := ports[port]; !ok {
			continue
		}
		if _, err := syscall.Getpeername(fd); err != nil {
			continue
		}
		_ = syscall.Shutdown(fd, syscall.SHUT_RDWR)
	}
}
