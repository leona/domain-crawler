.PHONY: default
default: build

build-small:
	@go build -ldflags="-s -w" -o ./bin/crawler ./main.go
	@upx --brute --best --lzma ./bin/crawler

build-linux:
	@GOOS=linux go build -ldflags="-s -w" -o ./bin/crawler ./src/main.go

build-mac:
	@GOOS=darwin go build -o ./bin/crawler ./main.go

build-all:
	@sudo apt-get install upx-ucl
	@GOOS=windows go build -ldflags="-s -w" -o ./bin/crawler ./main.go
	@upx --brute --best --lzma ./cli/cli.exe
	@GOOS=linux go build -ldflags="-s -w" -o ./bin/crawler ./main.go
	@GOOS=darwin go build -ldflags="-s -w" -o ./bin/crawler ./main.go
	
build-windows:
	@GOOS=windows go build ./main.go