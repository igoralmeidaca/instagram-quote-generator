# Use the official Golang image as a builder
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the worker binary
RUN go build -o /worker ./worker/main.go

# Use a minimal base image for execution
FROM gcr.io/distroless/base

# Copy the compiled binary directly (no need to chmod)
COPY --from=builder /worker /worker

# Run the worker
CMD ["/worker"]
