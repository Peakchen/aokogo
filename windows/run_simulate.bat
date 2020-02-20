@echo off
call make_simulate.bat
cd ..
cd bin
simulate.exe
cd ..
pause