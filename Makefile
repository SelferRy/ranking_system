.PHONY: deps, generate, mockgen, db-up, db-down, db-reset, db-shell, db-logs


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

db-up:
	docker compose up -d

db-down:
	docker compose down

db-reset:
	docker compose down -v
	docker compose up -d

db-shell:
	psql -h localhost -U user -d ranking_test

db-logs:
	docker compose logs -f postgres

migr-up:
	docker compose up -d
	goose up

migr-down:
	goose down

repo-test:
	docker compose up -d
	goose up
	cd internal/infra/adapters/repository/postgres/postgres_test && go test -v
	goose down
	docker compose down
