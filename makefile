start:
		go run ./cmd/api/main.go

migrate:
		go run ./pkg/db/migration.go
race:
		go run --race ./cmd/server.go