BINARY := drawio

.PHONY: linux
linux:
	mkdir -p build/linux
	GOOS=linux GOARCH=amd64 go build -o build/linux/$(BINARY) *.go

.PHONY: darwin
darwin:
	mkdir -p build/osx
	GOOS=darwin GOARCH=amd64 go build -o build/osx/$(BINARY) *.go

.PHONY: build
build:  linux darwin