#!/bin/bash
tts_tag=$(sed -n '1p' CHANGELOG.md)

echo "正在删除旧文件"
git rm ./release/*

echo "正在编译${tts_tag}_mac_amd64"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./release/${tts_tag}_mac_amd64
echo "正在编译${tts_tag}_windowns_amd64.exe"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./release/${tts_tag}_windowns_amd64.exe
echo "正在编译${tts_tag}_linux_amd64"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/${tts_tag}_linux_amd64
echo "正在编译${tts_tag}_linux_arm64"
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./release/${tts_tag}_linux_arm64
echo "正在编译${tts_tag}_linux_amd64"
go build -o ./release/${tts_tag}_termux_arm64

git add ./release/*
git commit -m "${tts_tag}"
git push origin master
git add ./*
git commit -m "${tts_tag}"
git push origin master



