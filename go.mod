module github.com/solo-io/reporting-client

go 1.12

require (
	github.com/gogo/protobuf v1.3.0 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/solo-io/go-utils v0.10.16
	golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc // indirect
	golang.org/x/net v0.0.0-20191003171128-d98b1b443823 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20191002091554-b397fe3ad8ed // indirect
	golang.org/x/tools v0.0.0-20191003162220-c56b4b191e2d // indirect
	google.golang.org/grpc v1.24.0
)

replace (
	github.com/Sirupsen/logrus v1.0.5 => github.com/sirupsen/logrus v1.0.5
	github.com/Sirupsen/logrus v1.3.0 => github.com/Sirupsen/logrus v1.0.6
	github.com/Sirupsen/logrus v1.4.2 => github.com/sirupsen/logrus v1.0.6
)
