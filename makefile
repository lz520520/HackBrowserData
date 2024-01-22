

build:
	go build -trimpath -o out/data.exe -ldflags "-w -s" cmd/bin/main.go


dll:
	GO_ENABLED=1 GOARCH=amd64 go build -trimpath -v  -ldflags="-s -w -extldflags=-Wl,e:/code/go/github/HackBrowserData/cmd/lib/dllmain.exp" -buildmode=c-shared -o ./out/data.dll  ./cmd/lib/