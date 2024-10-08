package main

//#include "run.h"
import "C"

import (
    "os"
    "tests/cmd/run"
    "tests/logger"
    "unsafe"
)

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
    logger.LogCallback = Log
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
            logger.LogErr(err.Error())
            return
        }
        run.Username = params.MustGetStringParam("username")
        run.MasterKey = params.MustGetStringParam("masterkey")
    }
    run.Execute()
}

func main() {
}
