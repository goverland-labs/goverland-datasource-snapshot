.PHONY: generate build

generate:
	protoc --proto_path=proto proto/*.proto --go_out=./proto --go-grpc_out=./proto

build:
	go mod tidy && go build && go test -v ./...
