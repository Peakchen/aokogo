@echo off
call make_GameServer.bat
cd ..
cd bin
GameServer.exe
cd ..
pause