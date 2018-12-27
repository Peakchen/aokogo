@echo off

set VERSION=v.0.1.0
echo version=%VERSION%

set GOPATH=D:\AokoPro\aoko
set GOPATH_BAK=%GOPATH%
set GOPATH=%GOPATH%;%~dp0;%~dp0\src\extra\;

set GOOS=windows
set GOARCH=amd64


echo start install sever ...
go install -gcflags " -N -l" GameServer

echo make ok

set GOPATH=%GOPATH_BAK%

pause