#!/bin/bash -x
goCmd=$(which go)

rm /bin/koffer 2>/dev/null
rm -rf /root/koffer 2>/dev/null
mkdir -p /tmp/bin

${goCmd} mod download
${goCmd} build -o bin/koffer

mv ./dev ./bin/koffer 2>/dev/null
cp -f ./bin/koffer /usr/bin/koffer 2>/dev/null
cp -f ./bin/koffer /tmp/bin/koffer 2>/dev/null
#exit_code="$?"
#echo "exit code: ${exit_code}"
#exit 0
