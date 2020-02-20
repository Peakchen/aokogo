@echo off
call make_GameServer.bat
cd ..
cd bin
start GameServer.exe
cd ..
pause