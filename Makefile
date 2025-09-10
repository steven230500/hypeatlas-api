dev:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

docker-build:
	docker build -t hypeatlas-api .
