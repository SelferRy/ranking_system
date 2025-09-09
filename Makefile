generate:
	protoc ./api/banner.proto --go_out=. --go-grpc_out=.