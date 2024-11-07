# Start with the official Go image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o app .

# Final stage
FROM alpine:latest

# Set working directory in the final container
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/app .

# Copy .env file into the container
COPY .env .env

# Make sure the .env file can be read by the application
RUN chmod 600 .env

# Expose port (modify according to your application)
EXPOSE 3000 
EXPOSE 3001 

# Command to run the application
CMD ["./app"]

