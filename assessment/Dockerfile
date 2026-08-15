FROM golang:1.26.4 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.22

WORKDIR /app

# Install certificates
RUN apk add --no-cache ca-certificates

# Copy binary
COPY --from=builder /app/server .

# Expose application port
EXPOSE 8080

# Start application
CMD ["./server"]
