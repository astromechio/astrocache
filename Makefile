
proto:
	protoc -I=server/master/service/proto --go_out=plugins=grpc:server/master/service service.proto
	protoc -I=server/verifier/service/proto --go_out=plugins=grpc:server/verifier/service service.proto
	protoc -I=server/worker/service/proto --go_out=plugins=grpc:server/worker/service service.proto