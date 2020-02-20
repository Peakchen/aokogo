@echo off

@echo "select redis 1"
redis-cli.exe select 1

@echo "clean data:"
redis-cli.exe flushdb
redis-cli.exe flushall

@echo "mongo oper clean£º"

mongo mongodb://127.0.0.1:27017/Server1 --eval "db.dropDatabase()"

pause