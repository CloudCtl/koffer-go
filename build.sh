#!/bin/bash
# cobra init --pkg-name github.com/CodeSparta/koffer-go
# cobra add mirror
# cobra add bundle
# go build
# gitup devel

goCmd=$(which go)

rm /bin/koffer 2>/dev/null
rm -rf /root/koffer 2>/dev/null
mkdir -p /tmp/bin

git stage -A; git commit -m "$@"; git push origin
${goCmd} get -u github.com/CodeSparta/koffer-go
${goCmd} get -u github.com/CodeSparta/koffer-go/cmd
${goCmd} build

mv ./dev koffer 2>/dev/null
cp -f koffer /usr/bin/koffer 2>/dev/null
mv -f koffer /tmp/bin/koffer 2>/dev/null

