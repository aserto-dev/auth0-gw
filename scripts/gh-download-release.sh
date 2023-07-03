#!/usr/bin/env bash

org=$1
repo=$2
goos=$3
goarch=$4

if [ $goos == "linux" ] && [ $goarch == "amd64" ] 
then
    curl -L -H "Accept: application/vnd.github+json" https://api.github.com/repos/$org/$repo/releases/latest | \
    jq '.assets[] | select(.name|test("linux_x86")).browser_download_url' | \
    xargs -I {} wget "{}" -O bin/ds-load-amd64.zip
elif [ $goos == "linux" ] && [ $goarch == "arm64" ]
then
    curl -L -H "Accept: application/vnd.github+json" https://api.github.com/repos/$org/$repo/releases/latest | \
    jq '.assets[] | select(.name|test("linux_arm64")).browser_download_url' | \
    xargs -I {} wget "{}" -O bin/ds-load-arm64.zip
fi
