module aoko

go 1.12

// github link latest
// for example: github.com/pkg/sftp latest
// go clean -modcache 清除缓存
// go mod vendor 自动创建vendor 目录

require (
	//golang.org/x/net/proxy latest
	github.com/Shopify/sarama v1.26.1
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394
	github.com/bsm/sarama-cluster v2.1.15+incompatible
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.18+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.3.3
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/gonutz/ide v0.0.0-20180502124734-e9fc8c14ed56
	github.com/google/btree v1.0.0 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/gorilla/websocket v1.4.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.13.0 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/klauspost/compress v1.10.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/prometheus/client_golang v1.4.1 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20190923202752-2cc03de413da
	github.com/satori/go.uuid v1.2.0
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/stvp/tempredis v0.0.0-20181119212430-b82af8480203 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200122045848-3419fae592fc // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	github.com/zieckey/etcdsync v0.0.0-20180810020013-cd5b26bc05a1
	go.etcd.io/bbolt v1.3.3 // indirect
	go.uber.org/zap v1.14.0 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/genproto v0.0.0-20200225123651-fc8f55426688 // indirect
	google.golang.org/grpc v1.27.1 // indirect
	gopkg.in/redsync.v1 v1.0.1
	sigs.k8s.io/yaml v1.2.0 // indirect
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
