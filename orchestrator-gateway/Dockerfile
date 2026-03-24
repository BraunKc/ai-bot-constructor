FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -trimpath -o ai-orchestrator-gateway ./cmd/main.go

# ---
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ai-orchestrator-gateway .
COPY /config/config.yml ./config/config.yml
COPY .env .

ENTRYPOINT ["./ai-orchestrator-gateway"]
