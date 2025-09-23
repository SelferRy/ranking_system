.PHONY: deps, generate, mockgen, db-up, db-down, db-reset, db-shell, db-logs


go_mod:
	go mod tidy
	go mod vendor

generate:
	protoc ./api/banner.proto --go_out=. --go-grpc_out=.

mockgen:
	mockgen -source=internal/domain/interfaces/repository/banner.go -destination=internal/mocks/mock_banner_repo.go -package=mocks
	mockgen -source=internal/domain/interfaces/repository/stats.go -destination=internal/mocks/mock_stats_repo.go -package=mocks
	mockgen -source=internal/domain/interfaces/broker/event_producer.go -destination=internal/mocks/mock_producer_interface.go -package=mocks
	mockgen -source=internal/domain/service/bandit/ucb1.go -destination=internal/mocks/mock_selector_service.go -package=mocks
	mockgen -source=internal/infra/adapters/broker/kafka/producer.go -destination=internal/mocks/mock_producer_implementation.go -package=mocks -mock_names=brokerWriter=MockBrokerWriter

db-test-up:
	docker compose -f tests/postgres-compose/docker-compose.yaml up -d

db-test-down:
	docker compose -f tests/postgres-compose/docker-compose.yaml down

db-test-reset:
	docker compose -f tests/postgres-compose/docker-compose.yaml down -v
	docker compose -f tests/postgres-compose/docker-compose.yaml up -d

db-test-shell:
	psql -h localhost -U user -d ranking_test

db-test-logs:
	docker compose -f tests/postgres-compose/docker-compose.yaml logs -f postgres

migr-test-up:
	docker compose -f tests/postgres-compose/docker-compose.yaml up -d
	goose up

migr-test-down:
	goose down

repo-test:
	docker compose -f tests/postgres-compose/docker-compose.yaml up -d
	goose up
	sleep 3s
	cd internal/infra/adapters/repository/postgres && go test -tags=integration -v
	goose down
	docker compose -f tests/postgres-compose/docker-compose.yaml down

kafka-test-up:
	docker compose -f tests/kafka-compose/docker-compose.yaml up -d

kafka-test-down:
	docker compose -f tests/kafka-compose/docker-compose.yaml down