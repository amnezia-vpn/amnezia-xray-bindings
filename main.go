package main

/*
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#ifndef AMNEZIA_XRAY_DEF
#define AMNEZIA_XRAY_DEF

#ifdef _MSC_VER
// MSVC does not support complex type definitions
// So here is a workaround to bypass this behavior
// See go.dev/issues/36233
typedef float _Fcomplex;
typedef double _Dcomplex;
#endif

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

	"github.com/xtls/xray-core/app/log"
	cLog "github.com/xtls/xray-core/common/log"
	"github.com/xtls/xray-core/transport/internet"

	"github.com/xtls/xray-core/core"
	_ "github.com/xtls/xray-core/main/distro/all"
)

var server *core.Instance = nil

//export amnezia_xray_free
func amnezia_xray_free(ptr unsafe.Pointer) {
	C.free(ptr)
}

//export amnezia_xray_configure
func amnezia_xray_configure(cConfig *C.char) *C.char {
	strConfig := C.GoString(cConfig)
	cfgReader := bytes.NewReader([]byte(strConfig))

	coreConfig, err := core.LoadConfig("json", cfgReader)
	if err != nil {
		return C.CString(err.Error())
	}

	server, err = core.New(coreConfig)
	if err != nil {
		return C.CString(err.Error())
	}

	return nil
}

//export amnezia_xray_start
func amnezia_xray_start() *C.char {
	if server == nil {
		return nil
	}

	if err := server.Start(); err != nil {
		return C.CString(err.Error())
	}

	return nil
}

//export amnezia_xray_stop
func amnezia_xray_stop() *C.char {
	if server == nil {
		return nil
	}

	if err := server.Close(); err != nil {
		return C.CString(err.Error())
	}

	return nil
}

//export amnezia_xray_setsockcallback
func amnezia_xray_setsockcallback(cb C.amnezia_xray_sockcallback, ctx unsafe.Pointer) *C.char {
	err := internet.RegisterDialerController(func(net, addr string, conn syscall.RawConn) error {
		conn.Control(func(fd uintptr) {
			C.amnezia_xray_invokesockcallback(cb, C.uintptr_t(fd), ctx)
		})
		return nil
	})

	if err != nil {
		return C.CString(err.Error())
	}

	return nil
}

func LogHandlerCreator(cb C.amnezia_xray_loghandler, ctx unsafe.Pointer) cLog.WriterCreator {
	return func() cLog.Writer {
		return &logHandler{cb: cb, ctx: ctx}
	}
}

type logHandler struct {
	cb  C.amnezia_xray_loghandler
	ctx unsafe.Pointer
}

func (l *logHandler) Write(msg string) error {
	cMsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cMsg))

	C.amnezia_xray_invokeloghandler(l.cb, cMsg, l.ctx)
	return nil
}

func (l *logHandler) Close() error {
	return nil
}

//export amnezia_xray_setloghandler
func amnezia_xray_setloghandler(cb C.amnezia_xray_loghandler, ctx unsafe.Pointer) {
	log.RegisterHandlerCreator(log.LogType_Console, func(lt log.LogType, options log.HandlerCreatorOptions) (cLog.Handler, error) {
		return cLog.NewLogger(LogHandlerCreator(cb, ctx)), nil
	})
}

func main() {
}
