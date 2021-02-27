# Start from the latest golang base image
FROM golang:alpine as builder

# Add Maintainer Info
LABEL maintainer="Alec Scott <alecbcs@github.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 go build -o binoc .

# Start again with minimal envoirnment.
FROM scratch

# Set the Current Working Directory inside the container
WORKDIR /app

COPY --from=builder /app/binoc /app/binoc

# Command to run the executable
ENTRYPOINT ["/app/binoc"]