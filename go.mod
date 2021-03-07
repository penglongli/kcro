module github.com/penglongli/kcro

go 1.14

require (
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.16.0
	google.golang.org/appengine v1.4.0
	google.golang.org/genproto v0.0.0-20210303154014-9728d6b83eeb // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.2.0
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
