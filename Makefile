run:
	go run cmd/company-service/main.go

build:
	go build -o bin/company-service ./cmd/company-service/main.go

docker-build:
	docker build -t company-service .

docker-up:
	docker compose up --build

test:
	go test ./...

swag:
	swag init --parseDependency --parseInternal -g ./cmd/company-service/main.go -o ./docs
