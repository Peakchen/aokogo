set DIR=%~dp0
cd /d "%DIR%"
cd message

echo pb make start...
setlocal enabledelayedexpansion
for %%i in (*.proto) do ( 
	  rem echo %%i 
	  set pbname=%DIR%go\ 
	  protoc -I %DIR%message\  --go_out=!pbname! %%i
)

pause