@echo off

echo "welcome use aoko!"

go fmt ./src/...

call .\\windows\\make_GameServer.bat

call .\\windows\\make_Externalgws.bat

call .\\windows\\make_Innergws.bat

call .\\windows\\make_loginserver.bat

call .\\windows\\make_simulate.bat

pause