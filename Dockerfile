# Build stage
FROM golang:1.24-alpine AS builder

# Set the working directory
WORKDIR /app

# Install CGO dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the project source code
COPY . .

# Build the application with CGO enabled
# -o /app/donatello: specifies that the compiled binary will be named 'donatello' and located in /app
# ./cmd/donatello: path to the main package for building
RUN CGO_ENABLED=1 go build -o /app/donatello ./cmd/donatello

# Runtime stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the build stage
COPY --from=builder /app/donatello .

# Copy the static resources required for the application
COPY resources ./resources

# Expose the port the service will run on
EXPOSE 8080

# Set the command to run the application
CMD ["./donatello"]