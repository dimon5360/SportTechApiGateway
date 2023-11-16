#!/bin/bash

unameOut=$(uname -a)
case "${unameOut}" in
    Darwin*)    protoc=protoc;;
    MINGW*)     protoc=D:/projects/protoc/bin/protoc.exe;;
    *Msys)      protoc=D:/projects/protoc/bin/protoc.exe;;
esac

${protoc} --version

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

export PATH="$PATH:$(go env GOPATH)/bin"

rm -rf proto/*.pb.go


${protoc} --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/*.proto
