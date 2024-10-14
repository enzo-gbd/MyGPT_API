FROM alpine:latest as perms
COPY /docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

FROM golang:1.23.2 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o GBA ./cmd/GBA/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migration ./cmd/migrate/migrate.go

FROM golang:1.23.2 as development

RUN go install github.com/air-verse/air@latest

COPY --from=builder /app /app

WORKDIR /app

COPY .air.toml .

CMD ["air"]

FROM alpine

COPY --from=perms /entrypoint.sh /entrypoint.sh
COPY --from=builder /app/GBA .
COPY --from=builder /app/migration .
COPY --from=builder /app/local.env ./app/
COPY --from=builder /app/tls/ ./tls/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
