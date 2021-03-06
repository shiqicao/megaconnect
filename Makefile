PROTOC = protoc

COVERAGE_OUT = coverage.out
# Keep adding to this list as we implement tests to target more packages
COVERPKG = ./chainmanager/...,./common/...,./connector/...,./flowmanager/...,./workflow/...,./unsafe/...

GOCC = gocc

PARSER_DEBUG = false

.PHONY: build install clean test cov covhtml dep protos cleanprotos parser cleanparser

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
	$(PROTOC) --go_out=paths=source_relative,plugins=grpc:. grpc/*.proto
	$(PROTOC) --go_out=paths=source_relative:. protos/*.proto

cleanprotos:
	find grpc protos -name '*.pb.go' -print -delete

parser:
	$(GOCC) -v -a -zip -debug_lexer=$(PARSER_DEBUG) -debug_parser=$(PARSER_DEBUG)  -o workflow/parser/gen -p github.com/megaspacelab/megaconnect/workflow/parser/gen workflow/parser/lang.bnf

cleanparser: 
	rm -rf workflow/parser/gen
