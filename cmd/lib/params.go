package main

import (
    "bytes"
    "encoding/binary"
)

func NewControlParams() *ControlParams {
    params := &ControlParams{
        meta: make(ParamMeta),
    }
    return params
}

type ParamMeta map[string]interface{}
type ControlParams struct {
    meta ParamMeta
}

func (this *ControlParams) FormatUnMarshal(data []byte) (err error) {
    readBuffer := bytes.NewBuffer(data)
    header := make([]byte, 4)
    readBuffer.Read(header)
    for {
        if readBuffer.Len() == 0 {
            break
        }
        // 获取key
        var keyLength byte
        keyLength, err = readBuffer.ReadByte()
        if err != nil {
            return
        }
        keyBytes := make([]byte, keyLength)
        _, err = readBuffer.Read(keyBytes)
        if err != nil {
            return
        }
        key := string(keyBytes)
        // 获取value type
        _, err = readBuffer.ReadByte()
        if err != nil {
            return
        }
        // 获取value length
        valueLengthBytes := make([]byte, 4)
        _, err = readBuffer.Read(valueLengthBytes)
        if err != nil {
            return
        }
        var valueLength uint32
        bytesBuffer := bytes.NewBuffer(valueLengthBytes)
        err = binary.Read(bytesBuffer, binary.LittleEndian, &valueLength)
        if err != nil {
            return
        }
        // 获取value
        valueBytes := make([]byte, valueLength)
        _, err = readBuffer.Read(valueBytes)
        if err != nil {
            return
        }
        this.meta[key] = valueBytes
    }
    return
}
func (this *ControlParams) MustGetStringParam(key string) (value string) {
    if v, ok := this.meta[key]; ok {
        switch vv := v.(type) {
        case []byte:
            value = string(vv)
        case string:
            value = vv
        }
    } else {
        value = ""
    }
    return
}
