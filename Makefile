PROTOC = protoc

COVERAGE_OUT = coverage.out
# Keep adding to this list as we implement tests to target more packages
COVERPKG = ./chainmanager/...,./common/...,./connector/...,./flowmanager/...,./workflow/...


.PHONY: build install clean test cov covhtml dep protos cleanprotos

build:
	go build ./...

install:
	go install ./...

clean:
	go clean ./...
	rm -f $(COVERAGE_OUT)

test:
	go test -race ./...

cov $(COVERAGE_OUT):
	go test -race -coverprofile=$(COVERAGE_OUT) -coverpkg=$(COVERPKG) ./...
	go tool cover -func=$(COVERAGE_OUT)

covhtml: $(COVERAGE_OUT)
	go tool cover -html=$<

dep:
	dep ensure

protos:
	$(PROTOC) -Igrpc --go_out=plugins=grpc:grpc grpc/*.proto

cleanprotos:
	rm -f grpc/*.pb.go .protos
