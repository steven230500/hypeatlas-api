db-up:
\tdocker-compose up -d

db-down:
\tdocker-compose down

migrate-up:
\t$(shell go env GOPATH)/bin/goose -dir ./db/migrations postgres "$(POSTGRES_URL)" up

migrate-down:
\t$(shell go env GOPATH)/bin/goose -dir ./db/migrations postgres "$(POSTGRES_URL)" down

run:
\tSTORAGE=postgres POSTGRES_URL=$(POSTGRES_URL) go run ./cmd/api

worker:
\tPOSTGRES_URL=$(POSTGRES_URL) go run ./cmd/worker

migrate-up:
\t$(shell go env GOPATH)/bin/goose -dir ./db/migrations postgres "$(POSTGRES_URL)" up

migrate-status:
\t$(shell go env GOPATH)/bin/goose -dir ./db/migrations postgres "$(POSTGRES_URL)" status
