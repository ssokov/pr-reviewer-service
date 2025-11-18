.PHONY: test test-unit test-integration test-coverage test-all lint fmt install-hooks

test:
	go test -short -v ./...

test-unit:
	go test -short -v ./internal/service/...

test-integration:
	docker-compose -f deployments/docker/docker-compose.yml up -d db
	@echo "Waiting for database..."
	@sleep 3
	TEST_DATABASE_URL="postgres://max_superuser:max_superuser@localhost:5433/pr_system?sslmode=disable" go test -v -coverprofile=coverage_integration.out -coverpkg=./... ./internal/repository/postgres

test-all:
	go test -v ./...

test-coverage:
	go test -short -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep total
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint:
	golangci-lint run --config .golangci.yml

fmt:
	go fmt ./...
	gci write $$(find . -name '*.go' -type f -not -path "*/vendor/*" -not -path "*/docs/*")

install-hooks:
	@echo "Installing git hooks..."
	@chmod +x deployments/git-hooks/pre-commit
	@cp deployments/git-hooks/pre-commit .git/hooks/pre-commit
	@echo "Pre-commit hook installed"

build:
	go build -o bin/pr-reviewer-service ./cmd/pr-reviewer-service

run:
	go run ./cmd/pr-reviewer-service

docker-up:
	docker-compose -f deployments/docker/docker-compose.yml up -d

docker-down:
	docker-compose -f deployments/docker/docker-compose.yml down

docker-build:
	docker build -f deployments/docker/Dockerfile -t pr-reviewer-service:latest .

migrate-up:
	migrate -path migrations -database "postgres://max_superuser:max_superuser@localhost:5433/pr_system?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://max_superuser:max_superuser@localhost:5433/pr_system?sslmode=disable" down