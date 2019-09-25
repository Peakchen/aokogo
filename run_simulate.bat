@echo off
call make_simulate.bat
cd bin
simulate.exe
cd ..
pause