@echo off
echo building generator
go build -ldflags="-s -w" -o osml.exe cmd/main.go
copy osml.exe c:\tools\