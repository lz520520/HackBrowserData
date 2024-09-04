package fileutil

import (
    "bytes"
    "io"
    "os"
    "path/filepath"
    "strings"
)

func FileContentSplit(content string) (resultSlice []string) {
    resultSlice = make([]string, 0)
    resultSlice = strings.Split(content, "\n")
    for i, c := range resultSlice {
        resultSlice[i] = strings.TrimSpace(c)
    }
    return
}

// 获取文件名，如C:/test.exe 返回 test
func GetBaseName(name string) string {
    filenameWithSuffix := filepath.Base(name)
    fileSuffix := filepath.Ext(filenameWithSuffix)
    return strings.TrimSuffix(filenameWithSuffix, fileSuffix)
}

// 获取文件名，如C:/test.exe 返回 test
func GetBaseNameWithSuffix(name string) string {
    filenameWithSuffix := filepath.Base(name)
    return filenameWithSuffix
}

func ReadFileBytes(filename string) (resultBytes []byte, err error) {
    resultBytes = make([]byte, 0)
    f, err := os.Open(filename)
    if err != nil {
        return
    }
    defer f.Close()
    buf := make([]byte, 1024)
    for {
        n, err := f.Read(buf)
        if err != nil && err != io.EOF {
            return resultBytes, err
        }
        if n == 0 {
            break
        }

        resultBytes = append(resultBytes, buf[:n]...)
    }
    return resultBytes, nil
}

func WriteFile(filename string, writeBytes []byte) (err error) {
    MkdirFromFile(filename)
    f, err := os.Create(filename)
    if err != nil {
        return
    }
    defer f.Close()
    f.Write(writeBytes)
    return nil
}

func AppendFile(filename string, writeBytes []byte) (err error) {
    var file *os.File
    file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    defer file.Close()
    _, err = file.Write(writeBytes)
    return
}
func CreateFile(filename string) (f *os.File, err error) {
    MkdirFromFile(filename)
    f, err = os.Create(filename)
    return

}
func MkdirFromFile(src string) error {
    dstDir := filepath.Dir(src)
    _, err := os.Stat(dstDir)
    if err != nil {
        err = os.MkdirAll(dstDir, os.ModePerm)
        if err != nil {
            return err
        }
    }
    return nil
}

func MoveFile(dstName, srcName string) error {
    // 尝试将源文件移动到目标位置
    err := os.Rename(srcName, dstName)
    if err != nil {
        return err
    }
    return nil
}

func CopyFileAttachBytes(dstName, srcName string, attach []byte) (err error) {
    input, err := os.ReadFile(srcName)
    if err != nil {
        return
    }
    input = append(input, attach...)

    err = os.WriteFile(dstName, input, 0644)
    return
}

func FileHasSuffix(filename string, subBytes []byte) bool {
    input, err := os.ReadFile(filename)
    if err != nil {
        return false
    }
    return bytes.HasSuffix(input, subBytes)

}

func GetAllFile(pathname string, suffix string) ([]string, error) {
    rd, err := os.ReadDir(pathname)
    if err != nil {
        return nil, err
    }
    result := make([]string, 0)
    for _, fi := range rd {
        if !fi.IsDir() {
            //fullName := pathname + "/" +
            if strings.HasSuffix(fi.Name(), suffix) {
                result = append(result, fi.Name())
            }
        }
    }
    return result, nil

}

func GetAllSubFile(pathname string, suffix []string) ([]string, error) {
    result := make([]string, 0)
    rd, err := os.ReadDir(pathname)
    if err != nil {
        return result, err
    }
    for _, fi := range rd {
        if fi.IsDir() {
            //fmt.Printf("[%s]\n", pathname+"\\"+fi.Name())
            tmpResult, err := GetAllSubFile(filepath.Join(pathname, fi.Name()), suffix)
            if err == nil {
                result = append(result, tmpResult...)
            }
        } else {
            filterFlags := false
            for _, ext := range suffix {
                if strings.TrimLeft(filepath.Ext(fi.Name()), ".") == ext {
                    filterFlags = true
                }
            }
            if !filterFlags {
                result = append(result, filepath.Join(pathname, fi.Name()))
            }
            //fmt.Println(filepath.Join(pathname,fi.Name()))
        }
    }
    return result, nil
}

func GetCurrentDir() string {
    dir, _ := os.Executable()
    exPath := filepath.Dir(dir)
    return exPath
}
