.PHONY: deps, generate

go_mod:
	go mod tidy
	go mod vendor

generate:
	protoc ./api/banner.proto --go_out=. --go-grpc_out=.

mockgen:
	mockgen -source=internal/domain/interfaces/repository/banner.go -destination=internal/mocks/mock_banner_repo.go -package=mocks