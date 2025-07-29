package main

import (
    "fmt"
    "github.com/syndtr/goleveldb/leveldb"
)

func main() {
    db, _ := leveldb.OpenFile("C:\\Users\\Administrator\\Downloads\\agent (95)\\c__users_25104_AppData_Local_local", nil)
    defer db.Close()
    iter := db.NewIterator(nil, nil)
    for iter.Next() {
        key := iter.Key()

        if string(key) == "_wocloud_storage_" {
            fmt.Println("key:", string(key))
        }
    }

}
