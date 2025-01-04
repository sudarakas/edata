build:
	@go build -o bin/edata cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/edata

migration:
	@migrate create -ext postgres -dir cmd/migrate/migrations $(filter-out $@, $(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Error: version is required. Usage: make migrate-force version=<version_number>"; \
		exit 1; \
	fi
	@go run cmd/migrate/main.go force $(version)

migrate-fix:
	@if [ -z "$(version)" ]; then \
		echo "Error: version is required. Usage: make migrate-fix version=<version_number>"; \
		exit 1; \
	fi
	@echo "Fixing dirty database state..."
	@go run cmd/migrate/main.go force $(version)
	@echo "Retrying migration down..."
	@go run cmd/migrate/main.go down
	@echo "Retrying migration up..."
	@go run cmd/migrate/main.go up