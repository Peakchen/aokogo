## go-torch -u  http://localhost:12004/debug/pprof/heap --colors mem --raw  -f mem.svg
## go-torch -u  http://localhost:12004 --seconds 60 --raw -f cpu.svg
## go tool pprof -raw -seconds 60 http://localhost:12004/debug/pprof/profile