set DIR=%~dp0
cd /d "%DIR%"
cd proto

echo pb make start...
setlocal enabledelayedexpansion
for %%i in (*.proto) do ( 
	  rem echo %%i 
	  set jsname=%DIR%js\ 
	  protoc -I %DIR%proto\  --js_out=!jsname! %%i
)

pause
echo "Done"
PAUSE