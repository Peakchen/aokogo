@echo off
call make_Innergws.bat
cd ..
cd bin
InnerGateway.exe
cd ..
pause