PROTOC = protoc
COVERAGE_OUT = coverage.out


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
	go test -race -coverprofile=$(COVERAGE_OUT) -covermode=atomic ./...

covhtml: $(COVERAGE_OUT)
	go tool cover -html=$<

dep:
	dep ensure

protos:
	$(PROTOC) -Igrpc --go_out=plugins=grpc:grpc grpc/*.proto

cleanprotos:
	rm -f grpc/*.pb.go .protos
