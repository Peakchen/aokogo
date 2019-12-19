@echo off

set VERSION=v.0.1.0
echo version=%VERSION%

set GOPATH_BAK=%GOPATH%
set GOPATH=%GOPATH%;%~dp0;%~dp0\aoko;

set GOOS=windows
set GOARCH=amd64

cd %~dp0src\vendor\
for /r %%i in (*.go) do (
	echo %%i
	go tool compile -I %~dp0\pkg\windows_amd64 %%i
)

echo make ok

set GOPATH=%GOPATH_BAK%

pause