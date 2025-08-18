FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o app cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app .

COPY --from=builder /app/internal/templates ./internal/templates
COPY --from=builder /app/static ./static

EXPOSE 8080
ENTRYPOINT ["./app"]
