#!/bin/bash

xc() {
	echo "$3_$1_$2"
	env GOOS=$1 GOARCH=$2 go build -o ./bin/$3_$1_$2$4 ./$3.go
}

echo "Cross build for MacOSX, Windows & Linux as 32bit & 64bit"
rm -rf bin/
xc darwin 386 tunnel
xc darwin amd64 tunnel
xc windows 386 tunnel .exe
xc windows amd64 tunnel .exe
xc linux 386 tunnel
xc linux amd64 tunnel
echo "Done!"