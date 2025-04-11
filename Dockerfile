# Use the official Go image as a base
FROM golang:1.24

# Install PostgreSQL client (psql) and other dependencies
RUN apt-get update && apt-get install -y postgresql-client

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application code
COPY . .

# Expose the port the app runs on
EXPOSE 8080

# Command to compile and run the application
CMD ["sh", "-c", "go build -o main ./cmd/main.go && ./main"]

ENTRYPOINT ["/app/scripts/deploy.sh"]
