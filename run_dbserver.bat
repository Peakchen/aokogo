@echo off
call make_dbserver.bat
cd bin
DBServer.exe
cd ..
pause