package main

//#include "run.h"
import "C"

import (
	"os"
	"tests/cmd/run"
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
	run.LogCallback = Log
}

//export DLLWMain
func DLLWMain(argsList uintptr) {
	//Log(0xf, "test", 4)
	run.Execute()
}

func main() {
}
