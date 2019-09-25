@echo off
call make_GameServer.bat
cd bin
GameServer.exe
cd ..
pause