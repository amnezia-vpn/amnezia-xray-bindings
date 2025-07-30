package main

/*
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#ifndef LIBXRAY_DEF
#define LIBXRAY_DEF

typedef void (*amnezia_xray_sockcallback)(uintptr_t fd, void* ctx);
typedef void (*amnezia_xray_loghandler)(char* str, void* ctx);

static inline void amnezia_xray_invokesockcallback(amnezia_xray_sockcallback cb, uintptr_t fd, void* ctx)
{
	cb(fd, ctx);
}

static inline void amnezia_xray_invokeloghandler(amnezia_xray_loghandler cb, char* msg, void* ctx)
{
	cb(msg, ctx);
}

#endif
*/
import "C"

import (
	"bytes"
	"syscall"
	"unsafe"

	"github.com/amnezia-vpn/amnezia-xray-core/common/log"
	"github.com/amnezia-vpn/amnezia-xray-core/transport/internet"

	"github.com/amnezia-vpn/amnezia-xray-core/core"
	_ "github.com/amnezia-vpn/amnezia-xray-core/main/distro/all"
)

const NoError = 0
const GenericError = -1

var server *core.Instance = nil
var lastError error

//export amnezia_xray_configure
func amnezia_xray_configure(cConfig *C.char) C.int {
	strConfig := C.GoString(cConfig)
	cfgReader := bytes.NewReader([]byte(strConfig))

	var coreConfig *core.Config
	coreConfig, lastError = core.LoadConfig("json", cfgReader)
	if lastError != nil {
		return GenericError
	}

	server, lastError = core.New(coreConfig)
	if lastError != nil {
		return GenericError
	}

	return NoError
}

//export amnezia_xray_start
func amnezia_xray_start() C.int {
	if server == nil {
		return GenericError
	}

	lastError = server.Start()
	if lastError != nil {
		return GenericError
	}

	return NoError
}

//export amnezia_xray_stop
func amnezia_xray_stop() C.int {
	if server == nil {
		return NoError
	}

	lastError = server.Close()
	if lastError != nil {
		return GenericError
	}

	return NoError
}

//export amnezia_xray_setsockcallback
func amnezia_xray_setsockcallback(cb C.amnezia_xray_sockcallback, ctx unsafe.Pointer) C.int {
	lastError = internet.RegisterDialerController(func(net, addr string, conn syscall.RawConn) error {
		conn.Control(func(fd uintptr) {
			C.amnezia_xray_invokesockcallback(cb, C.uintptr_t(fd), ctx)
		})
		return nil
	})

	if lastError != nil {
		return GenericError
	}

	return NoError
}

type logHandler struct {
	cb  C.amnezia_xray_loghandler
	ctx unsafe.Pointer
}

func (l *logHandler) Handle(msg log.Message) {
	cMsg := C.CString(msg.String())
	defer C.free(unsafe.Pointer(cMsg))

	C.amnezia_xray_invokeloghandler(l.cb, cMsg, l.ctx)
}

//export amnezia_xray_setloghandler
func amnezia_xray_setloghandler(cb C.amnezia_xray_loghandler, ctx unsafe.Pointer) {
	log.RegisterHandler(&logHandler{cb: cb, ctx: ctx})
}

func main() {
}
