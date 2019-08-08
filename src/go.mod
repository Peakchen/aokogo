module aoko

go 1.12

// github link latest
// for example: github.com/pkg/sftp latest
// go clean -modcache 清除缓存

require (
	github.com/acroca/go-symbols v0.1.1
	github.com/alecthomas/repr v0.0.0-20181024024818-d37bc2a10ba1
	github.com/antlinker/go-cmap v0.0.0-20160407022646-0c5e57012e96
	github.com/antlinker/go-dirtyfilter v1.2.0
	github.com/dutchcoders/goftp v0.0.0-20170301105846-ed59a591ce14
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	//github.com/go-delve/delve v1.2.0
	github.com/go-redsync/redsync v1.3.0
	//github.com/go latest
	github.com/go-zookeeper/zk v1.0.1
	github.com/golang/protobuf v1.3.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/gorilla/websocket v1.4.0
	github.com/kardianos/govendor v1.0.9
	github.com/karrick/godirwalk v1.10.12
	github.com/kr/fs v0.1.0
	github.com/mdempsky/gocode v0.0.0-20190203001940-7fb65232883f
	github.com/pkg/sftp v1.10.0
	github.com/ramya-rao-a/go-outline v0.0.0-20181122025142-7182a932836a
	github.com/redigo/redigo v0.0.0-20141115112439-201510e60683
	github.com/robfig/config v0.0.0-20141207224736-0f78529c8c7e
	github.com/rogpeppe/godef v1.1.1
	github.com/satori/go.uuid v1.2.0

	github.com/spf13/pflag v1.0.3
	github.com/sqs/goreturns v0.0.0-20181028201513-538ac6014518
	github.com/stathat/consistent v1.0.0
	github.com/typings/typings v2.1.1+incompatible
	github.com/uber-go/zap v1.10.0
	github.com/uudashr/gopkgs v2.0.1+incompatible
	github.com/yuin/gopher-lua v0.0.0-20190514113301-1cd887cd7036
	golang.org/x/arch v0.0.0

	golang.org/x/crypto v0.0.0
	golang.org/x/net v0.0.0
	golang.org/x/sync v0.0.0
	golang.org/x/sys v0.0.0
	golang.org/x/text v0.3.0
	golang.org/x/tools v0.0.0
	golang.org/x/xerrors v0.0.0
)

replace (
	golang.org/x/arch => github.com/golang/arch v0.0.0-20190312162104-788fe5ffcd8c
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/net => github.com/golang/net v0.0.0-20190724013045-ca1201d0de80
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190804053845-51ab0e2deafa
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190806215303-88ddfcebc769
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20190717185122-a985d3407aa7
)
