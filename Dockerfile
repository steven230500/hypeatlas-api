# syntax=docker/dockerfile:1.6

# --- build stage ---
FROM golang:1.25 AS build
WORKDIR /src

# Permite que Go descargue la toolchain que pida el go.mod si hace falta
ENV CGO_ENABLED=0 GO111MODULE=on GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Fuerza binario para linux/amd64 (tu droplet es x86_64)
ENV GOOS=linux GOARCH=amd64
RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg/mod \
  go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# --- runtime ---
FROM alpine:3.20
RUN addgroup -S app && adduser -S app -G app \
  && apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/api /app/api
USER app
EXPOSE 8080
CMD ["/app/api"]
