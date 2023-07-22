# Stage 1: Build the Go application
FROM golang:latest AS build

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . ./

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app .

# Stage 2: Create a lightweight container to run the application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go binary from the previous stage
COPY --from=build /go/src/github.com/rganes5/maanushi_earth_e-commerce/cmd/api/app ./

# Copy the .env file into the container
COPY .env ./

# Expose the application's port
EXPOSE 3000

# Run the Go application when the container starts
CMD ["./app"]
