# Build Stage
FROM golang:1.22 AS builder

WORKDIR /app

# Install ca-certificates in the build image (this will be available during build)
RUN apt-get update && apt-get install -y ca-certificates

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code files
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o app-binary ./cmd

# Final Stage (Minimal Image)
FROM scratch

# Copy CA certificates from the build image
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the build image
COPY --from=builder /app/app-binary .

# Expose the application port
EXPOSE 8088

# Run the application
CMD ["./app-binary"]
