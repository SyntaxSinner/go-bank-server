postgres:   # Run PostgreSQL container with specified user and password
	docker run --name postgres_db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:   # Create a database named 'simple_bank' within the PostgreSQL container
	docker exec -it postgres_db createdb --username=root --owner=root simple_bank

dropdb:     # Drop the 'simple_bank' database from the PostgreSQL container
	docker exec -it postgres_db dropdb simple_bank

migrate_up: # Apply database migrations to the 'simple_bank' database
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate_down:   # Roll back the last applied database migration for 'simple_bank'
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:   # Generate Go code using the SQLC tool
	sqlc generate

test:   # Run tests for the application
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrate_up migrate_down   # Declare specified targets as phony
