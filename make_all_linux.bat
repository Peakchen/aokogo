go fmt ./src/...

cd linux

call .\\make_dbserver_linux.bat

call .\\make_Externalgws_linux.bat

call .\\make_GameServer_linux.bat

call .\\make_Innergws_linux.bat

call .\\make_loginserver_linux.bat

call .\\make_simulate_linux.bat

cd ..

pause