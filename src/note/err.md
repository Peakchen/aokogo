- 1. can't load package: package xx: unknown import path ...
如果没有使用mod 则设置set GO111MODULE=auto, 如果有mod，则不能直接在main下直接编译，需要检查是否存在该导入包

- 2. runtime: out of memory: cannot allocate xxxx-byte block (xxx in use)
改用buff 相关内存管理，比如字符串使用可以bytes.Buffer大字符串处理