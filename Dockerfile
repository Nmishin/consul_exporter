# Stage 1: Build
FROM golang:1.22-alpine as builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o consul_exporter .

# Stage 2: Run
FROM busybox:1.35

WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/consul_exporter /app/consul_exporter

# Expose the port the application runs on
EXPOSE 9107

# Set the entrypoint to the Go binary
ENTRYPOINT ["./consul_exporter"]

