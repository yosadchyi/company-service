# Build stage
FROM golang:1.23-alpine AS builder

# Install git, curl, make
RUN apk add --no-cache git curl make

# Set working directory
WORKDIR /app

# Install swag CLI v1.8.12
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12

# Copy Go dependencies and Makefile
COPY go.mod go.sum Makefile ./
RUN go mod download

# Copy source
COPY . .

# Generate Swagger docs and build app
RUN make swag
RUN make build

# Final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/bin/company-service .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/configs ./configs

EXPOSE 8080

CMD ["./company-service"]