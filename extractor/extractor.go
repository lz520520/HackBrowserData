package extractor

import "tests/browserdata/master_keys"

// Extractor is an interface for extracting data from browser data files
type Extractor interface {
    Extract(masterKeys master_keys.MasterKeys) error

    Name() string

    Len() int
}
