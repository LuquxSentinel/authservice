build:
	@go build -o ./bin/authservice

run: build
	@./bin/authservice

test:
	@go test ./...