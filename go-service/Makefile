.PHONY: run test lint generate migrate-up migrate-down build docker

run:
	go run ./cmd/server

test:
	go test ./... -race -cover

lint:
	golangci-lint run ./...

# Regenerates internal/repo/sqlcgen from db/queries/*.sql + db/migrations/*.sql.
# Requires: https://docs.sqlc.dev/en/latest/overview/install.html
generate:
	sqlc generate

migrate-up:
	goose -dir db/migrations postgres "$$DATABASE_URL" up

migrate-down:
	goose -dir db/migrations postgres "$$DATABASE_URL" down

build:
	go build -o bin/server ./cmd/server

docker:
	docker build -t glovebox-go-service .
