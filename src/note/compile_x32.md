set path (for gcc)
set gcc 
set cc=xxx-gcc (win32) 
go mod
set GOOS=windows 
set GOARCH=386 
set GOHOSTARCH=386 
set GOPATH_BAK=%~dp0 
set GOPATH=%~dp0 
set gopath
set GO111MODULE
set CGO_ENABLED (if it about import `C`, then set 1.)

go build ...

focus：do not add " " in last command.