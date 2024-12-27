@ECHO off
echo %~dp0
echo %cd%
set GOARCH=amd64
go env -w GOARCH=amd64
#set GOOS=linux
#go env -w GOOS=linux
go build -o uf.exe
set GOOS=windows
go env -w GOOS=windows
@ECHO on