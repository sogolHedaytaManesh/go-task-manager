BENCH_PKG=./internal/service

###########################################
## Run Project in the docker environment ##
###########################################
.PHONY: docker-up docker-down
## Start all docker services
docker-up:
	docker-compose -f docker-compose.yml up -d

## Stop and remove all docker services
docker-down:
	docker-compose -f docker-compose.yml down --remove-orphans

.PHONY: build
## Build the project + postgres custom image without starting any services
build:
	docker-compose -f docker-compose.yml build

.PHONY: stop
## Stop and remove all services
stop:
	docker compose -f docker-compose.yml down --remove-orphans

.PHONY: restart
## Rebuild and Restart the go applications (master, minion, orchestrator)
restart:
	docker-compose -f docker-compose.yml up -d --build

.PHONY: doc
## Generate Swagger API documentation
doc:
	swag init -g internal/http/router.go --parseDependency --parseInternal --generatedTime=true

#############################################
## Performs tests in the local environment ##
#############################################

.PHONY: test
## Run the tests with a summary output, it's not supposed to be called directly
test:
	@go install gotest.tools/gotestsum@latest
	gotestsum --format testdox --format-icons hivis --format-hide-empty-pkg ${GOTEST_EXTRA_ARGS} $(or $(TEST_DIRECTORY),./...)

.PHONY: all-tests
## Run all unit tests
all-tests:
	go clean -testcache
	GOTEST_EXTRA_ARGS="-- -timeout=30m -count=1 -v" TEST_DIRECTORY="./..." make test

.PHONY: unit-tests
## Run all unit tests
unit-tests:
	go clean -testcache
	GOTEST_EXTRA_ARGS="-- -timeout=30m -count=1 -short -race" TEST_DIRECTORY="./..." make test

.PHONY: integration-tests
## Run all integrations tests
integration-tests:
	go clean -testcache
	GOTEST_EXTRA_ARGS="-- -timeout=30m -count=1 -run Integration -race -v" TEST_DIRECTORY="./..." make test

#############################################
## Performs pprof in the local environment ##
#############################################

benchmark-cpu:
	mkdir -p docs/profiling
	go test $(BENCH_PKG) -bench=. -benchmem -cpuprofile docs/profiling/cpu.prof

benchmark-heap:
	mkdir -p docs/profiling
	go test $(BENCH_PKG) -bench=. -benchmem -memprofile docs/profiling/heap.prof


# Run all benchmark tests
benchmarks:
	go clean -testcache
	go test -bench=. -benchmem -timeout=30m ./internal/service


###########################################
## Run migrations in the docker environment ##
###########################################

.PHONY: db-up
## Database: Bring up the Postgres database container
## Creates the necessary directories and starts the container
## Usage: make db-up
db-up:
	mkdir -p .docker/data/postgres; docker-compose -f docker-compose.yml up --build -d db

.PHONY: db-migrate-up
## Apply database migrations, useful for running the application without docker using go run ...
## Usage: make db-migrate-up
db-migrate-up: db-up
	docker run -v ${PWD}/database/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://myuser:mypass@localhost/mydb?sslmode=disable" up

.PHONY: db-migrate-down
## Rollback all database migrations,  useful for running the application without docker  go run ...
## Usage: make db-migrate-down
db-migrate-down:
	docker run -v ${PWD}/database/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://myuser:mypass@localhost/mydb?sslmode=disable" down -all

.PHONY: db-test-up
## Database TEST: Bring up the Postgres test database container
## Usage: make db-test-up
db-test-up:
	mkdir -p .docker/data/postgres_test; \
	docker-compose -f docker-compose.yml up --build -d db_test

.PHONY: db-test-migrate-up
## Apply database migrations to the test database
## Usage: make db-test-migrate-up
db-test-migrate-up: db-test-up
	docker run -v ${PWD}/database/migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ \
		-database "postgres://test_user:test_pass@localhost:6433/test_db?sslmode=disable" up

.PHONY: db-test-migrate-down
## Rollback all database migrations on test database
## Usage: make db-test-migrate-down
db-test-migrate-down:
	docker run -v ${PWD}/database/migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ \
		-database "postgres://test_user:test_pass@localhost:6433/test_db?sslmode=disable" down -all

