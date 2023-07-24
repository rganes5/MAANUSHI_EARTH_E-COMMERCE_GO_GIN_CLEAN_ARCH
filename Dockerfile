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
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main ./cmd/api

# Stage 2: Create a lightweight container to run the application
FROM gcr.io/distroless/base-debian11

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go binary from the previous stage
COPY --from=build /usr/src/app/main .

# Expose the application's port
EXPOSE 3000

# Run the Go application when the container starts
CMD ["./main"]
