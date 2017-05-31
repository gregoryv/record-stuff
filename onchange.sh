#!/bin/sh

path=$1
dir=$(dirname "$path")
filename=$(basename "$path")
extension="${filename##*.}"
nameonly="${filename%.*}"

case $extension in
    go)
        gofmt -w $path
	go-bindata -o data.go static/...
	go test -cover -coverprofile /tmp/c.out .
        go tool cover -html=/tmp/c.out -o /tmp/coverage.html
	go install github.com/gregoryv/record-stuff
        ;;
esac

