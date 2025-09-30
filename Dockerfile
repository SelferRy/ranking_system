FROM golang:1.24 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags='-s -w' -o /app/ranking_system ./cmd/ranking_system

FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata postgresql-client
COPY --from=builder /app/ranking_system /usr/local/bin/ranking_system
COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY configs /app/configs
COPY migrations /app/migrations
COPY seeds /app/seeds
COPY .env app/.
COPY docker/start.sh /app/start.sh
RUN chmod +x /app/start.sh

WORKDIR /app
EXPOSE 5080
ENTRYPOINT ["/app/start.sh"]
#ENTRYPOINT ["/usr/local/bin/ranking_system"]
#CMD ["serve"]