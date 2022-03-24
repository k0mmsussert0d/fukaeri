# build
FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./fukaeri ./cmd/fukaeri/main.go

ENTRYPOINT ["/app/fukaeri"]
