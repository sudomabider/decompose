#!/bin/sh

env GOOS=linux GOARCH=amd64 go build --ldflags="-s -w" -o dc.linux
env GOOS=darwin GOARCH=amd64 go build --ldflags="-s -w" -o dc.osx
env GOOS=windows GOARCH=amd64 go build --ldflags="-s -w" -o dc.exe
