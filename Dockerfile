# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o sentinel ./cmd/sentinel
# Runtime stage
FROM scratch

WORKDIR /

COPY --from=builder /app/sentinel /sentinel
COPY --from=builder /app/web /web

ENTRYPOINT ["/sentinel"]