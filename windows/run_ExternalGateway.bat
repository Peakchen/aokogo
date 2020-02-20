@echo off
call make_Externalgws.bat
cd ..
cd bin
ExternalGateway.exe
cd ..
pause