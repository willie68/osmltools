@echo off
echo building osml
goreleaser build --snapshot --clean --single-target
cd dist\osml_windows_amd64_v1
osml.exe version
rem go build -ldflags="-s -w" -o osml.exe cmd/osml/main.go
copy osml.exe c:\tools\
copy osml.exe ..\..
cd ..\..