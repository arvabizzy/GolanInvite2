.PHONY: run seed build test tidy

run:
	go run cmd/server/main.go

seed:
	go run cmd/seed/main.go

build:
	go build -v -o bin/server.exe cmd/server/main.go
	go build -v -o bin/seed.exe cmd/seed/main.go

test:
	go test -v ./...

tidy:
	go mod tidy
