GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

.PHONY: init
# init env
init:
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v3@v3.0.0-20260526000039-30da04b769dc
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v3@v3.0.0-20260526000039-30da04b769dc

GO_ERRORS_PROTO_PATH:=$(shell go env GOMODCACHE)/github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v3@v3.0.0-20260526000039-30da04b769dc

.PHONY: api
# generate api proto (kratos v3; errors must be v3 to avoid proto registration conflict)
api:
	cd ../.. && protoc --proto_path=./api \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:./api \
	       --go-grpc_out=paths=source_relative:./api \
	       --go-http_out=paths=source_relative:./api \
	       api/user-service/v1/user.proto
	cd ../.. && protoc --proto_path=./api \
	       --proto_path=$(GO_ERRORS_PROTO_PATH) \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:./api \
	       --go-errors_out=paths=source_relative:./api \
	       api/user-service/v1/errors.proto

.PHONY: config
# generate internal proto
config:
	buf generate --template buf.gen.config.yaml

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# generate
generate:
	go generate ./...
	go mod tidy

.PHONY: all
# generate all
all:
	make api
	make config
	make generate

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
