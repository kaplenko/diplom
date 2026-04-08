# ---- Stage 1: build the runner binary (stdlib only) ----
FROM golang:1.25-alpine AS builder

WORKDIR /build
RUN go mod init runner
COPY cmd/runner/main.go .
RUN CGO_ENABLED=0 go build -o /runner main.go

# ---- Stage 2: runtime with Go compiler for user code ----
FROM golang:1.25-alpine

RUN adduser -D -u 1000 sandbox
COPY --from=builder /runner /usr/local/bin/runner

USER sandbox
ENTRYPOINT ["runner"]
