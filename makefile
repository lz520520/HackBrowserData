



dll:
	CGO_ENABLED=1 GOARCH=amd64 go build -trimpath -v  -ldflags="-s -w -extldflags=-Wl,e:/code/go/github/HackBrowserData/cmd/lib/dllmain.exp" -buildmode=c-shared -o ./out/data_amd64.dll  ./cmd/lib/
	CGO_ENABLED=1 GOARCH=386 go build -trimpath -v  -ldflags="-s -w -extldflags=-Wl,e:/code/go/github/HackBrowserData/cmd/lib/dllmain_386.exp" -buildmode=c-shared -o ./out/data_386.dll  ./cmd/lib/