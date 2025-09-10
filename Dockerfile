FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api ./cmd/api

FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=build /app/api /app/
EXPOSE 8080
CMD ["/app/api"]
