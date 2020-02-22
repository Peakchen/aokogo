@echo off

set VERSION=v.0.1.0
echo version=%VERSION%

set GOPATH_BAK=%GOPATH%
set GOPATH=%GOPATH%;%~dp0;%~dp0..;

set GOOS=windows
set GOARCH=amd64


echo start install GameServer ...
go install -gcflags " -N -l" GameServer

rem cd %~dp0src\GameServer\

rem for /r %%i in (*.go) do (
rem	echo %%i
rem	go tool compile -I %~dp0\pkg\windows_amd64 %%i
rem )

rem go tool link -o ../../bin/GameServer.exe -L %~dp0\pkg\windows_amd64 main.o

echo make ok

set GOPATH=%GOPATH_BAK%