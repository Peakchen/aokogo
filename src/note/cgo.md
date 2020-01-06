#### cgo 编译过程
- 1. go tool cgo import_example.go 生成c的代码
- 2. 再调用gcc编译c代码
- 3. go tool cgo -dynimport 生成动态import
- 4. go tool compile -o example1.a -pack GO_FILES 生成动态库
- 5. go tool pack r example1.a _all.o 动态库打包
- 6. go tool link -o example example1.a 链接