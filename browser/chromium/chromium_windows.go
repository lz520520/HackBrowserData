//go:build windows

package chromium

import (
    "encoding/base64"
    "errors"
    "fmt"
    "os"
    "tests/logger"

    "github.com/tidwall/gjson"

    "tests/crypto"
    "tests/types"
    "tests/utils/fileutil"
)

var errDecodeMasterKeyFailed = errors.New("decode master key failed")

func (c *Chromium) GetMasterKey() ([]byte, error) {
    b, err := fileutil.ReadFile(types.ChromiumKey.TempFilename())
    if err != nil {
        return nil, err
    }
    defer os.Remove(types.ChromiumKey.TempFilename())

    encryptedKey := gjson.Get(b, "os_crypt.encrypted_key")
    if !encryptedKey.Exists() {
        return nil, nil
    }
    //logger.LogInfo(fmt.Sprintf("encrypt_key: %s", encryptedKey.String()))

    key, err := base64.StdEncoding.DecodeString(encryptedKey.String())
    if err != nil {
        return nil, errDecodeMasterKeyFailed
    }

    c.masterKey, err = crypto.DecryptWithDPAPI(key[5:])
    if err != nil {
        logger.Error("decrypt master key failed", "err", err)
        return nil, fmt.Errorf("decrypt master key failed, %v", err)
    }
    logger.Info("get master key success", "browser", c.name)
    return c.masterKey, nil
}
