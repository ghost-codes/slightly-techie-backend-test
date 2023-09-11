
test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: migrateup migratedown migrateup1 migratedown1 sqlc test server