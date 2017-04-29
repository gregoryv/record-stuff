#!/bin/sh

path=$1
dir=$(dirname "$path")
filename=$(basename "$path")
extension="${filename##*.}"
nameonly="${filename%.*}"

case $extension in
    go)
	gofmt -w $path
	go test -cover
        go install github.com/gregoryv/record-stuff
        ;;
esac

