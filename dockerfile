# Use the official Golang image as the base image
FROM golang:1.21.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Set the Git Access Key Arg
ARG ACCESS_KEY

# Install Git to clone the GitHub repository
RUN apk update && apk add git

# Clone your Go application repository
RUN git clone https://${ACCESS_KEY}@github.com/hngx-org/Demerzel-events-backend.git --depth=1 .

# Download modules
#RUN go mod download

# Build the app
RUN go build -o app

# Run the app
CMD ["./app"]