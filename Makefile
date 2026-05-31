build:
	@go build -o bin/rest-api-pg cmd/main.go

run: build
	@./bin/rest-api-pg

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@, $(MAKECMDGOALS))