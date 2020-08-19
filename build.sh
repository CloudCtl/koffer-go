#!/bin/bash -x
# cobra init --pkg-name github.com/CodeSparta/koffer-go
# cobra add mirror
# cobra add bundle
# go build
# gitup devel

goCmd=$(which go)

rm /bin/koffer 2>/dev/null
rm -rf /root/koffer 2>/dev/null
mkdir -p /tmp/bin

${goCmd} mod download

${goCmd} build

mv ./dev ./bin/koffer 2>/dev/null
cp -f ./bin/koffer /usr/bin/koffer 2>/dev/null
cp -f ./bin/koffer /tmp/bin/koffer 2>/dev/null

