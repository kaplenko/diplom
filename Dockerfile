# ---- Stage 1: Build ----
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/diplom ./cmd/diplom

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# ---- Target: migrate ----
FROM alpine:3.20 AS migrate

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY migrations /migrations

ENTRYPOINT ["goose", "-dir", "/migrations", "postgres"]

# ---- Target: runtime (default) ----
FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /bin/diplom .
COPY --from=builder /src/migrations ./migrations

EXPOSE 8080

CMD ["./diplom"]
