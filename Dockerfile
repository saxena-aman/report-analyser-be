# Use official Golang image for building the Go app
FROM golang:1.22-alpine AS build

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o report-analyser-be .

# Use a smaller image for running the built binary
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder image
COPY --from=build /app/report-analyser-be .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the application
CMD ["./report-analyser-be"]
