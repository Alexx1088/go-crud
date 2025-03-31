# Use the official Go image as a base
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the application code
COPY . .

# Download Go dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["/app/main"]
