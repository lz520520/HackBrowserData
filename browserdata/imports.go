// Package browserdata is responsible for initializing all the necessary
// components that handle different types of browser data extraction.
// This file, imports.go, is specifically used to import various data
// handler packages to ensure their initialization logic is executed.
// These imports are crucial as they trigger the `init()` functions
// within each package, which typically handle registration of their
// specific data handlers to a central registry.
package browserdata

import (
    _ "tests/browserdata/bookmark"
    _ "tests/browserdata/cookie"
    _ "tests/browserdata/creditcard"
    _ "tests/browserdata/download"
    _ "tests/browserdata/extension"
    _ "tests/browserdata/history"
    _ "tests/browserdata/localstorage"
    _ "tests/browserdata/password"
    _ "tests/browserdata/sessionstorage"
)
