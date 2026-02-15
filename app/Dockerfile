# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Install system packages
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary (CGO_ENABLED=0 makes it static for alpine/scratch)
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./main.go

# Stage 2: Final image
FROM alpine:3.20

WORKDIR /root/

# Copy executable from builder
COPY --from=builder /main .

# Expose gRPC port
EXPOSE 50051

CMD ["./main"]