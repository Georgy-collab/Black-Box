# Build Go binary
FROM golang:1.22-bookworm AS builder

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/user-service ./cmd/server

# Runtime: Go binary + Python validator
FROM python:3.12-slim-bookworm

WORKDIR /app

COPY --from=builder /out/user-service /app/user-service
COPY python ./python
COPY openapi.json openapi.yaml ./

ENV PORT=8000
ENV PYTHON_BIN=python3
ENV PYTHON_VALIDATOR=/app/python/validate_cli.py

EXPOSE 8000

CMD ["/app/user-service"]
