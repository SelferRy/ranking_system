.PHONY: deps, generate

go_mod:
	go mod tidy
	go mod vendor

generate:
	protoc ./api/banner.proto --go_out=. --go-grpc_out=.

mockgen:
	mockgen -source=internal/domain/interfaces/repository/banner.go -destination=internal/mocks/mock_banner_repo.go -package=mocks
	mockgen -source=internal/domain/interfaces/repository/stats.go -destination=internal/mocks/mock_stats_repo.go -package=mocks
	mockgen -source=internal/domain/interfaces/gateway/event_producer.go -destination=internal/mocks/mock_producer_gateway.go -package=mocks
	mockgen -source=internal/domain/service/bandit/ucb1.go -destination=internal/mocks/mock_selector_service.go -package=mocks
