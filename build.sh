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

plugins="
    "github.com/spf13/cobra" \
    "github.com/go-git/go-git" \
    "github.com/go-git/go-git/plumbing" \
    "github.com/CodeSparta/koffer-go/err" \
    "github.com/CodeSparta/koffer-go/log" \
    "github.com/CodeSparta/koffer-go/auth" \
    "github.com/CodeSparta/koffer-go/status" \
"
for i in ${plugins}; do
  ${goCmd} get -u ${i};
done

${goCmd} build

mv ./dev koffer 2>/dev/null
cp -f koffer /usr/bin/koffer 2>/dev/null
mv -f koffer /tmp/bin/koffer 2>/dev/null

