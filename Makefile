build:
	go build -o bin/main cmd/main/main.go

dev:
	GIN_MODE=debug go run cmd/main/main.go

test:
	gotestdox ./...