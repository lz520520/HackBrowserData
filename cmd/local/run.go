package main

import (
    "archive/zip"
    "bytes"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "tests/browser"
    "tests/browserdata/master_keys"
    "tests/log"
    "tests/utils/fileutil"

    "github.com/gogf/gf/v2/encoding/gcompress"
    "github.com/google/uuid"
    "github.com/urfave/cli/v2"
)

var (
    outputDir    string
    outputFormat string
    verbose      bool
    compress     bool
    profilePath  string
    isFullExport bool
    dataPath     string
)

type BrowserKey struct {
    Home       string `json:"home"`
    Browser    string `json:"browser"`
    DefaultKey []byte `json:"default_key"`
    V20Key     []byte `json:"v20_key"`
}

func Execute() {
    app := &cli.App{
        Name:      "tools",
        Usage:     "Export passwords|bookmarks|cookies|history|credit cards|download history|localStorage|extensions from browser",
        UsageText: "[tools -b chrome -f json -dir results --zip]\nExport all browsing data (passwords/cookies/history/bookmarks) from browser",
        Version:   "0.5.0",
        Flags: []cli.Flag{
            &cli.BoolFlag{Name: "verbose", Aliases: []string{"vv"}, Destination: &verbose, Value: false, Usage: "verbose"},
            &cli.BoolFlag{Name: "compress", Aliases: []string{"zip"}, Destination: &compress, Value: false, Usage: "compress result to zip"},
            &cli.StringFlag{Name: "results-dir", Aliases: []string{"dir"}, Destination: &outputDir, Value: "tmpout", Usage: "export dir"},
            &cli.StringFlag{Name: "format", Aliases: []string{"f"}, Destination: &outputFormat, Value: "json", Usage: "output format: csv|json"},
            &cli.StringFlag{Name: "profile-path", Aliases: []string{"p"}, Destination: &profilePath, Value: "", Usage: "custom profile dir path, get with chrome://version"},
            &cli.BoolFlag{Name: "full-export", Aliases: []string{"full"}, Destination: &isFullExport, Value: true, Usage: "is export full browsing data"},
            &cli.StringFlag{Name: "data", Destination: &dataPath, Value: "", Usage: "data path"},
        },
        HideHelpCommand: true,
        Action: func(c *cli.Context) error {
            tmpout, _ := filepath.Abs(filepath.Join(outputDir, uuid.New().String()))
            defer os.RemoveAll(tmpout)
            err := gcompress.UnZipFile(dataPath, tmpout)
            if err != nil {
                return err
            }
            keyInfos := make([]BrowserKey, 0)
            keyBytes, err := fileutil.ReadFileBytes(filepath.Join(tmpout, "key.json"))
            if err != nil {
                return err
            }
            err = json.Unmarshal(keyBytes, &keyInfos)
            if err != nil {
                return err
            }
            home := filepath.Join(tmpout, "Users")

            log.LogInfo("start collect data")
            if verbose {
                log.SetVerbose()
                //logger.Default.SetVerbose()
                //logger.Configure(logger.Default)
            }

            outBuffer := bytes.Buffer{}
            zw := zip.NewWriter(&outBuffer)

            for _, info := range keyInfos {
                userPath := filepath.Join(home, info.Home)
                browser.RefreshConfig(userPath)
                browsers, err := browser.PickBrowsers(info.Browser, "")
                if err != nil {
                    log.LogErr(fmt.Sprintf("pick [%s] browsers error: %s", info.Home, err.Error()))
                    log.Errorf("pick [%s] browsers %v", info.Home, err)
                    continue
                } else {
                    log.LogSuccess(fmt.Sprintf("[%s] PickBrowsers success", info.Home))
                }

                for _, b := range browsers {
                    log.LogSuccess(fmt.Sprintf("get [%s] browsing %s data", info.Home, b.Name()))
                    data, err := b.BrowsingData(isFullExport, info.Home, &master_keys.MasterKeys{
                        DefaultKey: info.DefaultKey,
                        V20Key:     info.V20Key,
                    })
                    if err != nil {
                        log.LogErr(fmt.Sprintf("get [%s] browsing data error: %s", info.Home, err.Error()))
                        log.Errorf("get [%s] browsing data error %v", info.Home, err)
                        continue
                    }
                    data.Output(zw, info.Home, b.Name(), outputFormat)
                }
            }
            zw.Close()
            log.LogInfo("over collect data")

            b := outBuffer.Bytes()
            log.LogBytes(b, len(b))
            resultPath := filepath.Join(outputDir, uuid.New().String()+".zip")
            err = fileutil.WriteFile(resultPath, b)
            if err != nil {
                return err
            }
            log.LogSuccess(fmt.Sprintf("result save: %s", resultPath))
            return nil
        },
    }
    err := app.Run(os.Args)
    if err != nil {
        log.LogErr(fmt.Sprintf("run app error %v", err))
    }
}
