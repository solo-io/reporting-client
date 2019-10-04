.PHONY: update-deps
update-deps:
	go get github.com/golang/protobuf/protoc-gen-go

SUBDIRS:=api
.PHONY: generated-code
generated-code:
	protoc api/v1/reporting.proto --go_out=plugins=grpc:./pkg
	gofmt -w $(SUBDIRS)
	goimports -w $(SUBDIRS)
