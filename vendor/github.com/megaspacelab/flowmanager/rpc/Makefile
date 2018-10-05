PROTOC = protoc


.PHONY: protos clean

protos: .protos

.protos: *.proto
	$(PROTOC) --go_out=plugins=grpc:. $^
	touch .protos

clean:
	rm -f *.pb.go .protos
