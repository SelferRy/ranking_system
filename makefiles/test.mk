.PHONY: db-test-up, db-test-down, db-test-reset, db-test-shell, db-test-logs

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

migr-test-up: db-test-up
	source .env.test && goose up
	sleep 3s

migr-test-down:
	source .env.test && goose down
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
