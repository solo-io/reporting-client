.PHONY: init
init:
	git config core.hooksPath .githooks

.PHONY: update-deps
update-deps:
	go get github.com/golang/protobuf/protoc-gen-go

SUBDIRS:=api
.PHONY: generated-code
generated-code:
	protoc api/v1/reporting.proto --go_out=plugins=grpc:./pkg
	go generate ./...
	gofmt -w $(SUBDIRS)
	goimports -w $(SUBDIRS)
