.PHONY: all build

ver := $(shell git log -1 --pretty=format:"%h-%as")

build:
	go build -ldflags "-w -s -X main.gitCommit=$(ver)" -o build/twc cli/main.go

all:
	GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X main.gitCommit=$(ver)" -o build/twc-linux-amd64 cli/main.go
	GOOS=linux GOARCH=arm64 go build -ldflags "-w -s -X main.gitCommit=$(ver)" -o build/twc-linux-arm64 cli/main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "-w -s -X main.gitCommit=$(ver)" -o build/twc-mac-arm64 cli/main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s -X main.gitCommit=$(ver)" -o build/twc-mac-amd64 cli/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags "-w -s -X main.gitCommit=$(ver)" -o build/twc-x64.exe cli/main.go