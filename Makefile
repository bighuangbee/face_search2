GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

INTERNAL_PROTO_FILES=$(shell find pkg -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)

.PHONY: init
# init env
init:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2
	go get -u github.com/google/wire/cmd/wire
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go install github.com/envoyproxy/protoc-gen-validate@latest


.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./pkg \
		--proto_path=./third_party \
		--go_out=paths=source_relative:./pkg \
		$(INTERNAL_PROTO_FILES)


errors:
	protoc --proto_path=. \
	   --proto_path=../../../third_party \
	   --proto_path=../../../pkg/proto \
	   --go_out=paths=source_relative:. \
	   --go-errors_out=paths=source_relative:. \
	   $(API_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api \
	   --proto_path=./third_party \
	   --go_out=paths=source_relative:./api \
	   --go-http_out=paths=source_relative:./api \
	   --go-grpc_out=paths=source_relative:./api \
	   --validate_out=paths=source_relative,lang=go:./api \
	   --go-errors_out=paths=source_relative:./api \
	   --openapiv2_out=./api \
	   $(API_PROTO_FILES)

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...


.PHONY: test
# test
test:
	go test -v ./... -cover

.PHONY: run
run:
	cd app/cmd/registe/ && go build -o registeBin *.go
	cd app/cmd/server/ && go run .

.PHONY: docker
docker:
	docker build -f app/Dockerfile -t biz_svc .

.PHONY: wire
# generate wire
wire:
	cd app/cmd/server && wire

gorm_gen:
	rm -rf app/internal/data/dal && cd app/cmd/ && go run .

.PHONY: all
# generate all
all:
	make api;
	make config;
	make generate;

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

clean:
	docker rmi -f  `docker images | grep '<none>' | awk '{print $3}'`
