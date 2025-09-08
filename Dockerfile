# stage 1: compile and build the Go app
FROM golang:1.24.6-alpine AS builder


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env .env

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server ./cmd/server/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/seeder ./cmd/seed/main.go

#stage 2: run the container
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/seeder .

COPY --from=builder /app/sql ./sql
COPY --from=builder /app/.env .env

EXPOSE 8484

