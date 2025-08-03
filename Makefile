build:
	go build -o bin/main cmd/main/main.go

dev:
	go run cmd/main/main.go

test:
	go test -json -v $$(go list ./... | grep -E '/internal/|/pkg/') | gotestfmt