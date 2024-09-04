package main

import (
    "fmt"
    "github.com/syndtr/goleveldb/leveldb"
)

func main() {
    db, _ := leveldb.OpenFile("C:\\Users\\Administrator\\Downloads\\Procdump\\leveldb", nil)
    defer db.Close()
    iter := db.NewIterator(nil, nil)
    for iter.Next() {
        key := iter.Key()

        if string(key) == "_wocloud_storage_" {
            fmt.Println("key:", string(key))
        }
    }

}
