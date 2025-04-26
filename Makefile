migrate:
	go run cmd/main.go --with-migrate

run:
	go run cmd/main.go

clean:
	go clean -modcache
	go mod tidy
	go mod download

test:
	go test -cover ./...
