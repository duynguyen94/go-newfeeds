# Start from the official Golang base image
FROM golang:1.22-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
# Build the user Go service
RUN go build -o users-service ./cmd/users

# Build the newsfeed Go service
RUN go build -o newsfeed-service ./cmd/newsfeed

# Copy the entrypoint script into the container
COPY docker/entrypoint.sh /app/entrypoint.sh

# Make the entrypoint script executable
RUN chmod +x /app/entrypoint.sh

# Expose port 8080 to the outside world
EXPOSE 8080

# Set the custom entrypoint script
ENTRYPOINT ["/app/entrypoint.sh"]

# Optional: You can add CMD to pass default arguments to your entrypoint
# CMD ["--help"]  # Example default argument