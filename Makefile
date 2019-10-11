VERSION_INJECT=main.versionText
SRCS=*.go cmd/*.go
export GO111MODULE=on

all: linux-amd64 windows-amd64 darwin-amd64

linux-amd64: bin/opiekey_linux_amd64

windows-amd64: bin/opiekey_windows_amd64.exe

darwin-amd64: bin/opiekey_darwin_amd64

test:
	go test -coverprofile=coverage.txt -covermode=atomic . ./cmd
	go tool cover -html=coverage.txt -o coverage.html

bin/opiekey_windows_amd64.exe: $(SRCS)
	GOOS=windows GOARCH=amd64 go build -o $@ -ldflags "-X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)" github.com/arcanericky/opiekey/cmd

bin/opiekey_linux_amd64: $(SRCS)
	GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "-X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)" github.com/arcanericky/opiekey/cmd

bin/opiekey_darwin_amd64: $(SRCS)
	GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags "-X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)" github.com/arcanericky/opiekey/cmd

clean:
	rm -rf bin
