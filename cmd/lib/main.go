package main

//#include "run.h"
import "C"

import (
    "os"
    "runtime"
    _ "runtime"
    "runtime/debug"
    "tests/cmd/run"
    "tests/log"
    "time"
    "unsafe"
)

//go:linkname runtime_beforeExit os.runtime_beforeExit
func runtime_beforeExit(exitCode int)

func Log(logType int, msg []byte, size int) {
    cLogType := C.int(logType)
    cSize := C.int(size)
    if size > 0 {
        cMsg := C.CBytes(msg)
        defer C.free(cMsg)
        C.cPrint(cLogType, (*C.char)(cMsg), cSize)
    } else {
        cMsg := C.CString(string(msg))
        C.cPrint(cLogType, cMsg, cSize)
    }
}

//export install
func install(installArgs uintptr) {
    C.initArgs(unsafe.Pointer(installArgs))
    os.Args = []string{os.Args[0]}
    log.LogCallback = Log
}

//export DLLWMain
func DLLWMain(argsList uintptr) {
    if argsList > 0 {
        args := (*C.struct__args_list)(unsafe.Pointer(argsList))

        //Log(0xf, "test", 4)
        paramBytes := C.GoBytes(unsafe.Pointer((*args).params), (*args).params_length)
        params := NewControlParams()
        err := params.FormatUnMarshal(paramBytes)
        if err != nil {
            log.LogErr(err.Error())
            return
        }
        run.Username = params.MustGetStringParam("username")
        run.MasterKey = params.MustGetStringParam("masterkey")
        run.LibProfilePath = params.MustGetStringParam("profile_path")
        run.LibBrowserName = params.MustGetStringParam("browser_name")
    }
    run.Execute()
    clean()
}

func clean() {
    runtime.GC()
    debug.FreeOSMemory()
    time.Sleep(time.Second * 3)
    runtime_beforeExit(0)
}

func main() {
}
