# --- Build Stage ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build Master
RUN CGO_ENABLED=0 GOOS=linux go build -o master ./cmd/master

# Build Agent
RUN CGO_ENABLED=0 GOOS=linux go build -o agent ./cmd/agent

# --- Sentinel Master Image ---
FROM alpine:latest AS master
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/master .
COPY --from=builder /app/web ./web
EXPOSE 8080
ENTRYPOINT ["./master"]

# --- Sentinel Agent Image ---
FROM alpine:latest AS agent
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/agent .
ENTRYPOINT ["./agent"]
