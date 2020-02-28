@echo off

set VERSION=v.0.1.0
echo version=%VERSION%

set GOPATH_BAK=%GOPATH%
set GOPATH=%GOPATH%;%~dp0;%~dp0..;

set GOOS=linux
set GOARCH=amd64


echo start install GameServer ...
go install -gcflags " -N -l" GameServer

echo make ok

set GOPATH=%GOPATH_BAK%