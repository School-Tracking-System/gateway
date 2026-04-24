# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy the modules needed
COPY proto/ ./proto/
COPY services/gateway/ ./services/gateway/

# Create a service-specific go.work to avoid loading other services
RUN printf "go 1.25.0\n\nuse (\n\t./proto\n\t./services/gateway\n)\n" > go.work

# Build the application
WORKDIR /app/services/gateway
# Gateway might not have module.go if it follows a different structure, let's check
RUN go build -o bin/api cmd/api/main.go cmd/api/module.go

# Final stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/services/gateway/bin/api .
COPY --from=builder /app/services/gateway/.env.template .env

EXPOSE 8000

CMD ["./api"]
