run:
	go run cmd/main.go
build:
	CGO_ENABLED=0 go build -o server cmd/main.go
test: test-registry test-api
test-registry:
	go test -v ./internal/registry
test-api: test-storages test-goods
test-storages:
	go test -v ./internal/handler/storages
test-goods:
	go test -v ./internal/handler/goods
coverage:
	go test -v -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

up:
	docker compose up -d
down:
	docker compose down
destroy:
	docker compose down -v
restart:
	docker compose restart
