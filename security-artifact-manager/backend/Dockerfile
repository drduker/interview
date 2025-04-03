# Stage 1: Build the Go application
FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o artifact-manager .

# Stage 2: Create the final image
FROM alpine:3.17

# Add security updates and CA certificates
RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

# Set up non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Create directory for artifact storage with proper permissions
RUN mkdir -p /artifacts
VOLUME /artifacts

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/artifact-manager .

# Expose the application port
EXPOSE 8080

# Command to run the executable
CMD ["./artifact-manager"]