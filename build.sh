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
platforms="windows_amd64 linux_amd64 linux_arm"
for platform in $platforms; do

    rm -rf ./resources/output/rtmp2flv_"${ver}"_linux_amd64

    export GOOS=$(echo "$platform" | gawk 'BEGIN{FS="_"} {print $1}')
    export GOARCH=$(echo "$platform" | gawk 'BEGIN{FS="_"} {print $2}')
    export CGO_ENABLED=0
    echo "${GOOS}"_"${GOARCH}"
    if [[ "${GOOS}" == "windows" ]]
    then
        go build -o ./resources/output/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/rtmp2flv.exe main.go
    else
        go build -o ./resources/output/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/rtmp2flv main.go
    fi
    go build -o ./resources/output/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/rtmp2flv main.go

    mkdir -p ./resources/output/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/output/live
    mkdir -p ./resources/output/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/output/log
    mkdir -p ./resources/output/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/conf

    cp -r ./resources/static ./resources/resources/output/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/static/
    cp -r ./resources/db ./resources/output/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/db/
    cp -r ./resources/conf ./resources/output/rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/conf

    cd ./resources/output/ || exit
    rm -rf rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}".tar.gz
    tar -zcvf ./rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}".tar.gz ./rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/

    rm -rf ./rtmp2flv_"${ver}"_"${GOOS}"_"${GOARCH}"/
    cd ../
done
