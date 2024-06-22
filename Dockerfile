# Start from the official Go image
FROM golang:alpine

# Install bash and wget
RUN apk add --no-cache bash wget

# Download wait-for-it script
RUN wget -O /wait-for-it.sh https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh \
    && chmod +x /wait-for-it.sh

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main ./cmd/hex

# Expose the port the app runs on
EXPOSE 8080

# Use wait-for-it to wait for MySQL before starting the app
CMD ["/bin/bash", "-c", "/wait-for-it.sh mysql:3307 -- ./main"]