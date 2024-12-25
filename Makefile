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