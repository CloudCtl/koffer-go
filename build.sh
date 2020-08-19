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

#git stage -A; git commit -m 'testing'; git push origin master

plugins="
    "github.com/spf13/cobra" \
    "golang.org/x/sys/unix" \
    "github.com/spf13/viper" \
    "github.com/go-git/go-git" \
    "github.com/go-git/go-git/plumbing" \
    "github.com/CodeSparta/koffer-go/plugins/err" \
    "github.com/CodeSparta/koffer-go/plugins/log" \
    "github.com/CodeSparta/koffer-go/plugins/auth" \
"
for i in ${plugins}; do
  ${goCmd} get -u ${i};
done

${goCmd} build

mv ./dev ./bin/koffer 2>/dev/null
cp -f ./bin/koffer /usr/bin/koffer 2>/dev/null
cp -f ./bin/koffer /tmp/bin/koffer 2>/dev/null

