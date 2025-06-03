package main

import "tests/browserdata/master_keys"

func main() {
    master_keys.GetMasterKey("help_root", `C:\Users\Administrator\AppData\Local\Microsoft\Edge\User Data\Local State`)
}
