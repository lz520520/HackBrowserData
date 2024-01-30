package main

//#include "run.h"
import "C"

import (
	"os"
	"tests/cmd/run"
	"unsafe"
)

func Log(logType int, msg string, size int) {
	cLogType := C.int(logType)
	cMsg := C.CString(msg)
	cSize := C.int(size)
	C.cPrint(cLogType, cMsg, cSize)
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
