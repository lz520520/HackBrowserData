package log

import (
    "os"
    "path/filepath"
    "tests/utils/fileutil"
)

var (
    LogCallback func(logType int, msg []byte, size int) = Log
)

func Log(logType int, msg []byte, size int) {
    switch logType {
    case 0:
        Info(string(msg))
    case 1:
        Info(string(msg))
    case 2:
        Error(string(msg))
    case 0xf:
        name := filepath.Join(filepath.Dir(os.Args[0]), "data.bin")
        fileutil.WriteFile(name, msg)
    }
}

func LogSuccess(msg string) {
    if LogCallback != nil {
        LogCallback(1, []byte(msg), 0)
    }
}
func LogInfo(msg string) {
    if LogCallback != nil {
        LogCallback(0, []byte(msg), 0)
    }
}
func LogErr(msg string) {
    if LogCallback != nil {
        LogCallback(2, []byte(msg), 0)
    }
}
func LogBytes(msg []byte, size int) {
    if LogCallback != nil {
        LogCallback(0xf, msg, size)
    }
}
