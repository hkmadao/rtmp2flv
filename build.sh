#/bin/bash
ver=$1
if [ -n "${ver}" ] 
then
    echo package version ${ver}
else
    echo no version param
    exit 1
fi
# Linux
# export GOOS=linux
# export GOARCH=amd64
export CGO_ENABLED=1
go build -o rtmp2flv_${ver}_linux_amd64 main.go

# Windows
# export GOOS=windows
# export GOARCH=amd64
# export CGO_ENABLED=1
# go build -o rtmp2flv.exe main.go

#package linux_amd64
rm -rf ./output/rtmp2flv_${ver}_linux_amd64

mkdir -p ./output/rtmp2flv_${ver}_linux_amd64/output/live
mkdir -p ./output/rtmp2flv_${ver}_linux_amd64/output/log
mkdir -p ./output/rtmp2flv_${ver}_linux_amd64/conf

cp -r ./static ./output/rtmp2flv_${ver}_linux_amd64/static/
cp -r ./db ./output/rtmp2flv_${ver}_linux_amd64/db/
cp -r ./conf/conf.yml ./output/rtmp2flv_${ver}_linux_amd64/conf
cp -r ./rtmp2flv_${ver}_linux_amd64 ./output/rtmp2flv_${ver}_linux_amd64/rtmp2flv

#package window_amd64
rm -rf ./output/rtmp2flv_${ver}_window_amd64

mkdir -p ./output/rtmp2flv_${ver}_window_amd64/output/live
mkdir -p ./output/rtmp2flv_${ver}_window_amd64/output/log
mkdir -p ./output/rtmp2flv_${ver}_window_amd64/conf

cp -r ./static ./output/rtmp2flv_${ver}_window_amd64/static/
cp -r ./db ./output/rtmp2flv_${ver}_window_amd64/db/
cp -r ./conf/conf.yml ./output/rtmp2flv_${ver}_window_amd64/conf
cp -r ./rtmp2flv_${ver}_window_amd64.exe ./output/rtmp2flv_${ver}_window_amd64/rtmp2flv.exe
cp -r ./start.vbs ./output/rtmp2flv_${ver}_window_amd64/start.vbs

cd ./output/
rm -rf rtmp2flv_*.tar.gz
tar -zcvf ./rtmp2flv_${ver}_linux_amd64.tar.gz ./rtmp2flv_${ver}_linux_amd64/
tar -zcvf ./rtmp2flv_${ver}_window_amd64.tar.gz ./rtmp2flv_${ver}_window_amd64/

rm -rf ./rtmp2flv_${ver}_linux_amd64/
rm -rf ./rtmp2flv_${ver}_window_amd64/