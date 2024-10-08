package browser

import (
    "path/filepath"
    "sort"
    "strings"
    "tests/browserdata"
    "tests/logger"

    "tests/browser/chromium"
    "tests/browser/firefox"
    "tests/utils/fileutil"
    "tests/utils/typeutil"
)

type Browser interface {
    // Name is browser's name
    Name() string
    // BrowsingData returns all browsing data in the browser.
    BrowsingData(isFullExport bool, masterkey []byte) (*browserdata.BrowserData, error)
}

// PickBrowsers returns a list of browsers that match the name and profile.
func PickBrowsers(name, profile string) ([]Browser, error) {
    var browsers []Browser
    clist := pickChromium(name, profile)
    for _, b := range clist {
        if b != nil {
            browsers = append(browsers, b)
        }
    }
    flist := pickFirefox(name, profile)
    for _, b := range flist {
        if b != nil {
            browsers = append(browsers, b)
        }
    }
    return browsers, nil
}

func pickChromium(name, profile string) []Browser {
    var browsers []Browser
    name = strings.ToLower(name)
    if name == "all" {
        for _, v := range chromiumList {
            if !fileutil.IsDirExists(filepath.Clean(v.profilePath)) {
                logger.Warn("find browser failed, profile folder does not exist", "browser", v.name)
                continue
            }
            multiChromium, err := chromium.New(v.name, v.storage, v.profilePath, v.dataTypes)
            if err != nil {
                logger.Error("new chromium error", "err", err)
                continue
            }
            for _, b := range multiChromium {
                logger.Warn("find browser success", "browser", b.Name())
                browsers = append(browsers, b)
            }
        }
    }
    if c, ok := chromiumList[name]; ok {
        if profile == "" {
            profile = c.profilePath
        }
        if !fileutil.IsDirExists(filepath.Clean(profile)) {
            logger.Error("find browser failed, profile folder does not exist", "browser", c.name)
        }
        chromiumList, err := chromium.New(c.name, c.storage, profile, c.dataTypes)
        if err != nil {
            logger.Error("new chromium error", "err", err)
        }
        for _, b := range chromiumList {
            logger.Warn("find browser success", "browser", b.Name())
            browsers = append(browsers, b)
        }
    }
    return browsers
}

func pickFirefox(name, profile string) []Browser {
    var browsers []Browser
    name = strings.ToLower(name)
    if name == "all" || name == "firefox" {
        for _, v := range firefoxList {
            if profile == "" {
                profile = v.profilePath
            } else {
                profile = fileutil.ParentDir(profile)
            }

            if !fileutil.IsDirExists(filepath.Clean(profile)) {
                logger.Warn("find browser failed, profile folder does not exist", "browser", v.name)
                continue
            }

            if multiFirefox, err := firefox.New(profile, v.dataTypes); err == nil {
                for _, b := range multiFirefox {
                    logger.Warn("find browser success", "browser", b.Name())
                    browsers = append(browsers, b)
                }
            } else {
                logger.Error("new firefox error", "err", err)
            }
        }

        return browsers
    }

    return nil
}

func ListBrowsers() []string {
    var l []string
    l = append(l, typeutil.Keys(chromiumList)...)
    l = append(l, typeutil.Keys(firefoxList)...)
    sort.Strings(l)
    return l
}

func Names() string {
    return strings.Join(ListBrowsers(), "|")
}
