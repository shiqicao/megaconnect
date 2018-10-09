PROTOC = protoc


.PHONY: all build install clean dep protos cleanprotos

all: build

build:
	go build ./...

install:
	go install ./...

clean:
	go clean ./...

dep:
	dep ensure

protos:
	$(PROTOC) -Igrpc --go_out=plugins=grpc:grpc grpc/*.proto

cleanprotos:
	rm -f grpc/*.pb.go .protos
