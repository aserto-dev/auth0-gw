#!/usr/bin/env bash

if [ $1 == "linux" ] && [ $2 == "amd64" ] 
then
    cp bin/$2/ds-load dist/auth0-gw_linux_amd64_v1
    cp bin/$2/ds-load-auth0 dist/auth0-gw_linux_amd64_v1
elif [ $1 == "linux" ] && [ $2 == "arm64" ]
then
    cp bin/$2/ds-load dist/auth0-gw_linux_arm64
    cp bin/$2/ds-load-auth0 dist/auth0-gw_linux_arm64
else
    echo "DO NOTHING"
fi
