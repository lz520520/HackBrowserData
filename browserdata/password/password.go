package password

import (
    "bytes"
    "database/sql"
    "encoding/base64"
    "os"
    "sort"
    "tests/browserdata/master_keys"
    "time"

    "github.com/tidwall/gjson"
    _ "modernc.org/sqlite" // import sqlite3 driver

    "tests/crypto"
    "tests/extractor"
    "tests/log"
    "tests/types"
    "tests/utils/typeutil"
)

func init() {
    extractor.RegisterExtractor(types.ChromiumPassword, func() extractor.Extractor {
        return new(ChromiumPassword)
    })
    extractor.RegisterExtractor(types.YandexPassword, func() extractor.Extractor {
        return new(YandexPassword)
    })
    extractor.RegisterExtractor(types.FirefoxPassword, func() extractor.Extractor {
        return new(FirefoxPassword)
    })
}

type ChromiumPassword []loginData

type loginData struct {
    UserName    string
    encryptPass []byte
    encryptUser []byte
    Password    string
    LoginURL    string
    CreateDate  time.Time
}

const (
    queryChromiumLogin = `SELECT origin_url, username_value, password_value, date_created FROM logins`
)

func (c *ChromiumPassword) Extract(masterKeys master_keys.MasterKeys) error {
    db, err := sql.Open("sqlite", types.ChromiumPassword.TempFilename())
    if err != nil {
        return err
    }
    defer os.Remove(types.ChromiumPassword.TempFilename())
    defer db.Close()

    rows, err := db.Query(queryChromiumLogin)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var (
            url, username string
            pwd, password []byte
            create        int64
        )
        if err := rows.Scan(&url, &username, &pwd, &create); err != nil {
            log.Debugf("scan chromium password error: %v", err)
        }
        login := loginData{
            UserName:    username,
            encryptPass: pwd,
            LoginURL:    url,
        }

        if len(pwd) > 0 {
            if bytes.HasPrefix(pwd, []byte("v20")) {
                password, err = crypto.DecryptWithChromium(masterKeys.V20Key, pwd)
            } else {
                password, err = crypto.DecryptWithChromium(masterKeys.DefaultKey, pwd)
            }
            if err != nil {
                log.Debugf("decrypt chromium password error: %v", err)
            }
            if bytes.HasPrefix(pwd, []byte("v20")) && len(password) > 32 {
                password = password[32:]
            }
        }

        if create > time.Now().Unix() {
            login.CreateDate = typeutil.TimeEpoch(create)
        } else {
            login.CreateDate = typeutil.TimeStamp(create)
        }
        login.Password = string(password)
        *c = append(*c, login)
    }
    // sort with create date
    sort.Slice(*c, func(i, j int) bool {
        return (*c)[i].CreateDate.After((*c)[j].CreateDate)
    })
    return nil
}

func (c *ChromiumPassword) Name() string {
    return "password"
}

func (c *ChromiumPassword) Len() int {
    return len(*c)
}

type YandexPassword []loginData

const (
    queryYandexLogin = `SELECT action_url, username_value, password_value, date_created FROM logins`
)

func (c *YandexPassword) Extract(masterKeys master_keys.MasterKeys) error {
    db, err := sql.Open("sqlite", types.YandexPassword.TempFilename())
    if err != nil {
        return err
    }
    defer os.Remove(types.YandexPassword.TempFilename())
    defer db.Close()

    rows, err := db.Query(queryYandexLogin)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var (
            url, username string
            pwd, password []byte
            create        int64
        )
        if err := rows.Scan(&url, &username, &pwd, &create); err != nil {
            log.Debugf("scan yandex password error: %v", err)
        }
        login := loginData{
            UserName:    username,
            encryptPass: pwd,
            LoginURL:    url,
        }

        if len(pwd) > 0 {
            if bytes.HasPrefix(pwd, []byte("v20")) {
                password, err = crypto.DecryptWithChromium(masterKeys.V20Key, pwd)
            } else {
                password, err = crypto.DecryptWithChromium(masterKeys.DefaultKey, pwd)
            }
            if err != nil {
                log.Debugf("decrypt yandex password error: %v", err)
            }
            if bytes.HasPrefix(pwd, []byte("v20")) && len(password) > 32 {
                password = password[32:]
            }
        }
        if create > time.Now().Unix() {
            login.CreateDate = typeutil.TimeEpoch(create)
        } else {
            login.CreateDate = typeutil.TimeStamp(create)
        }
        login.Password = string(password)
        *c = append(*c, login)
    }
    // sort with create date
    sort.Slice(*c, func(i, j int) bool {
        return (*c)[i].CreateDate.After((*c)[j].CreateDate)
    })
    return nil
}

func (c *YandexPassword) Name() string {
    return "password"
}

func (c *YandexPassword) Len() int {
    return len(*c)
}

type FirefoxPassword []loginData

func (f *FirefoxPassword) Extract(masterKeys master_keys.MasterKeys) error {
    logins, err := getFirefoxLoginData()
    if err != nil {
        return err
    }

    for _, v := range logins {
        userPBE, err := crypto.NewASN1PBE(v.encryptUser)
        if err != nil {
            return err
        }
        pwdPBE, err := crypto.NewASN1PBE(v.encryptPass)
        if err != nil {
            return err
        }
        user, err := userPBE.Decrypt(masterKeys.DefaultKey)
        pwd, err := pwdPBE.Decrypt(masterKeys.DefaultKey)
        if err != nil {
            return err
        }
        *f = append(*f, loginData{
            LoginURL:   v.LoginURL,
            UserName:   string(user),
            Password:   string(pwd),
            CreateDate: v.CreateDate,
        })
    }

    sort.Slice(*f, func(i, j int) bool {
        return (*f)[i].CreateDate.After((*f)[j].CreateDate)
    })
    return nil
}

func getFirefoxLoginData() ([]loginData, error) {
    s, err := os.ReadFile(types.FirefoxPassword.TempFilename())
    if err != nil {
        return nil, err
    }
    defer os.Remove(types.FirefoxPassword.TempFilename())
    loginsJSON := gjson.GetBytes(s, "logins")
    var logins []loginData
    if loginsJSON.Exists() {
        for _, v := range loginsJSON.Array() {
            var (
                m    loginData
                user []byte
                pass []byte
            )
            m.LoginURL = v.Get("formSubmitURL").String()
            user, err = base64.StdEncoding.DecodeString(v.Get("encryptedUsername").String())
            if err != nil {
                return nil, err
            }
            pass, err = base64.StdEncoding.DecodeString(v.Get("encryptedPassword").String())
            if err != nil {
                return nil, err
            }
            m.encryptUser = user
            m.encryptPass = pass
            m.CreateDate = typeutil.TimeStamp(v.Get("timeCreated").Int() / 1000)
            logins = append(logins, m)
        }
    }
    return logins, nil
}

func (f *FirefoxPassword) Name() string {
    return "password"
}

func (f *FirefoxPassword) Len() int {
    return len(*f)
}
