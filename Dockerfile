# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.14-alpine as build

# Set the Current Working Directory inside the container
WORKDIR /src

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o /app/main.sh ./cmd
RUN rm -R /src
WORKDIR /app

# Expose port 9090 to the outside world
EXPOSE 9090

# Command to run the executable
CMD ["./main.sh"]