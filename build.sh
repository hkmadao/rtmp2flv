#!/bin/bash
#./build.sh 0.0.1
ver=$1
if [ -n "${ver}" ]; then
    echo package version "${ver}"
else
    echo no version param
    exit 1
fi
#打多个平台的包
platforms="linux_amd64"
rm -rf ./resources/output/releases/
for platform in $platforms; do

    export GOOS=$(echo "$platform" | gawk 'BEGIN{FS="_"} {print $1}')
    export GOARCH=$(echo "$platform" | gawk 'BEGIN{FS="_"} {print $2}')
    export CGO_ENABLED=1
    echo "${GOOS}"_"${GOARCH}"
    if [[ "${GOOS}" == "windows" ]]; then
        go build -o ./resources/output/releases/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/rtmp2flv.exe main.go
    else
        go build -o ./resources/output/releases/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/rtmp2flv main.go
    fi
    go build -o ./resources/output/releases/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/rtmp2flv main.go

    mkdir -p ./resources/output/releases/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/resources/output/live
    mkdir -p ./resources/output/releases/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/resources/output/log
    mkdir -p ./resources/output/releases/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/resources/conf

    cp -r ./resources/static ./resources/output/releases/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/resources/static/
    cp -r ./resources/db ./resources/output/releases/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/resources/db
    cp -r ./resources/conf ./resources/output/releases/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/resources/conf

    cd ./resources/output/releases/ || exit
    rm -rf rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}".tar.gz
    tar -zcvf ./rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}".tar.gz rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/

    # rm -rf ./rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/
    cd ../../../
done
