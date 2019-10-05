@echo off
call make_loginserver.bat
cd bin
LoginServer.exe
cd ..
pause