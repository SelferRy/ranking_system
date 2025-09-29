.PHONY: deps, generate, mockgen, db-up, db-down, db-reset, db-shell, db-logs

include makefiles/test.mk

build:
	go build ./cmd/ranking_system

run: cli-up

test:
	go test ./...

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

cli-compose-up:
	docker compose -f docker-compose.cli.yaml up -d
	sleep 3s

cli-compose-down:
	docker compose -f docker-compose.cli.yaml down -v

migr-up:
	goose up

seed:
	@echo ">>> Loading seed data..."
	docker exec -i $$(docker compose -f "docker-compose.cli.yaml" ps -q postgres) \
		psql -U prod_user -d ranking_system < seeds/cli_seed.sql

cli-deps-up: cli-compose-up migr-up seed

cli-deps-down: cli-compose-down

cli-up:
	$(MAKE) cli-deps-up
	sleep 4s
	go run ./cmd/ranking_system serve

cli-pipeline:
	$(MAKE) cli-data-deps-up
	go run ./cmd/ranking_system serve
	# check it in another terminal via:
	# grpcurl -plaintext -proto api/proto/banner_rotator.proto -d '{"slot_id": 1, "group_id": 1}' localhost:5080 ranking_system.BannerRotatorService/SelectBanner

db-migr-up:
	$(MAKE) cli-deps-up
	goose up
	sleep 3s

db-migr-down:
	goose down
	$(MAKE) cli-deps-down


# Remember:
# cli-up - запускает серверную часть
# проверка серверной части через: grpcurl -plaintext -proto api/proto/banner_rotator.proto -d '{"slot_id": 1, "group_id": 1}' localhost:5080 ranking_system.BannerRotatorService/SelectBanner
# пока не работает создание топика, надо внутри контейнера выполнить: kafka-topics --create --if-not-exists --bootstrap-server kafka:9092 --topic ranking_system_topic --partitions 1 --replication-factor 1
# и проверка: kafka-topics --bootstrap-server localhost:9092 --list
