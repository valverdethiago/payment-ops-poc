migrate-up:
	cd src;\
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/bank_transactional_db?sslmode=disable" -verbose up

migrate-down:
	cd src;\
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/bank_transactional_db?sslmode=disable" -verbose down

sqlc:
	cd src;\
	sqlc generate

test:
	cd src;\
	go test -v -cover ./...

build:
	cd src;\
	go build -o ../docker/bin/trio-provider-ms

run:
	cd src;\
	go run main.go