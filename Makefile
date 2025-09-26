.PHONY: deps, generate, mockgen, db-up, db-down, db-reset, db-shell, db-logs


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

db-test-up:
	@echo "Start database container"
	docker compose -f tests/postgres-compose/docker-compose.yaml up -d
	sleep 3s

db-test-down:
	@echo "Down database container"
	docker compose -f tests/postgres-compose/docker-compose.yaml down

db-test-reset:
	@echo "Restart database container with delete data volume"
	docker compose -f tests/postgres-compose/docker-compose.yaml down -v
	docker compose -f tests/postgres-compose/docker-compose.yaml up -d

db-test-shell:
	psql -h localhost -U user -d ranking_test

db-test-logs:
	docker compose -f tests/postgres-compose/docker-compose.yaml logs -f postgres

migr-test-up:
	$(MAKE) db-test-up
	goose up
	sleep 3s

migr-test-down:
	goose down
	$(MAKE) db-test-down

repo-test:
	$(MAKE) migr-test-up
	cd internal/infra/adapters/repository/postgres && go test -tags=integration -v
	$(MAKE) migr-test-down

kafka-test-up:
	docker compose -f tests/kafka-compose/docker-compose.yaml up -d
	sleep 3s

kafka-test-down:
	@echo "Down kafka container"
	docker compose -f tests/kafka-compose/docker-compose.yaml down

kafka-test:
	$(MAKE) kafka-test-up
	cd internal/infra/adapters/broker/kafka && go test -tags=integration -v
	$(MAKE) kafka-test-down

usecase-integration-test:
	$(MAKE) migr-test-up
	$(MAKE) kafka-test-up
	cd internal/domain/usecase/banner && go test -tags=integration -v
	$(MAKE) kafka-test-down
	$(MAKE) migr-test-down
