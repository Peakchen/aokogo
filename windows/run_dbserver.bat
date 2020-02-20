@echo off
call make_dbserver.bat
cd ..
cd bin
DBServer.exe
cd ..
pause