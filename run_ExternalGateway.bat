@echo off
call make_Externalgws.bat

cd bin
ExternalGateway.exe
cd ..
pause