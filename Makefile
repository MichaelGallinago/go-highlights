.PHONY: gen-client
	protoc --go_out=./mainservice/ --go-grpc_out=require_unimplemented_servers=false:./mainservice/  ./productservice/api/proto/api.proto


.PHONY: gen-server
	protoc --go_out=./productservice/ --go-grpc_out=./productservice/  ../productservice/api/proto/api.proto
