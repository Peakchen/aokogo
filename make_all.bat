@echo off

echo "welcome use aoko!"

call make_Externalgws.bat

call make_GameServer.bat

call make_Innergws.bat

call make_loginserver.bat

call make_simulate.bat
pause