build:
	@go build -o bin/go-api cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/go-api

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(NAME)

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

migrate-version:
	@go run cmd/migrate/main.go version

migrate-reset:
	@go run cmd/migrate/main.go -1 force

.PHONY: build test run migration migrate-up migrate-down migrate-version migrate-reset