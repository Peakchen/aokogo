@echo off

call make_Externalgws.bat
pause
call make_GameServer.bat
pause
call make_Innergws.bat
pause
call make_loginserver.bat
pause
call make_simulate.bat
pause