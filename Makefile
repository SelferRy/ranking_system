.PHONY: deps, generate, mockgen, db-up, db-down, db-reset, db-shell, db-logs

include makefiles/test.mk

go_mod:
	go mod tidy
	go mod vendor

generate:
	@echo "Run proto generate"
	protoc \
		--proto_path=api/proto \
		--go_out=api/gen --go_opt=paths=source_relative \
		--go-grpc_out=api/gen --go-grpc_opt=paths=source_relative \
		 ./api/proto/banner_rotator.proto

mockgen:
	mockgen -source=internal/domain/interfaces/repository/banner.go -destination=internal/mocks/mock_banner_repo.go -package=mocks
	mockgen -source=internal/domain/interfaces/repository/stats.go -destination=internal/mocks/mock_stats_repo.go -package=mocks
	mockgen -source=internal/domain/interfaces/broker/event_producer.go -destination=internal/mocks/mock_producer_interface.go -package=mocks
	mockgen -source=internal/domain/service/bandit/ucb1.go -destination=internal/mocks/mock_selector_service.go -package=mocks
	mockgen -source=internal/infra/adapters/broker/kafka/producer.go -destination=internal/mocks/mock_producer_implementation.go -package=mocks -mock_names=brokerWriter=MockBrokerWriter

migr-up:
	$(MAKE) db-test-up
	goose up
	sleep 3s

migr-down:
	goose down
	$(MAKE) db-test-down

