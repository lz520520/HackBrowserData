package master_keys

/*
#cgo LDFLAGS: -L. -lplugin_chrome -lntdll -loleaut32 -lbcrypt -lcrypt32 -lwtsapi32 -lncrypt -luserenv
#include "run.h"
#include <stdlib.h>
*/
import "C"
import (
    "fmt"
    "unsafe"
)

type MasterKeys struct {
    DefaultKey []byte
    V20Key     []byte
}

func ChromeInstall(args uintptr) {
    cArgs := (*C.install_args_t)(unsafe.Pointer(args))
    C.chrome_install(cArgs)
}

func GetMasterKey(username string, keyPath string) (MasterKeys, error) {
    masterKeys := MasterKeys{
        DefaultKey: make([]byte, 0),
        V20Key:     make([]byte, 0),
    }

    args := C.master_key_args_t{
        username: C.CString(username),
        key_path: C.CString(keyPath),
    }
    defer C.free(unsafe.Pointer(args.username))
    defer C.free(unsafe.Pointer(args.key_path))

    cArgs := (*C.master_key_args_t)(C.malloc(C.size_t(unsafe.Sizeof(args))))
    defer C.free(unsafe.Pointer(cArgs))
    *cArgs = args

    cStatus := C.chrome_get_master_key(cArgs)
    status := int(cStatus)
    if status != 0 {
        return masterKeys, fmt.Errorf("get master key error: %d\n", status)
    }
    if int((*cArgs).v10_key_len) > 0 {
        masterKeys.DefaultKey = C.GoBytes(unsafe.Pointer((*cArgs).v10_key), (*cArgs).v10_key_len)
    }
    if int((*cArgs).v20_key_len) > 0 {
        masterKeys.V20Key = C.GoBytes(unsafe.Pointer((*cArgs).v20_key), (*cArgs).v20_key_len)

    }
    return masterKeys, nil
}
