@echo off
echo building osml
go build -ldflags="-s -w" -o osml.exe cmd/osml/main.go
copy osml.exe c:\tools\