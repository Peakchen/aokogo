@echo off
call make_loginserver.bat
cd ..
cd bin
LoginServer.exe
cd ..
pause