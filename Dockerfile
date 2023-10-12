#FROM ubuntu:latest
#LABEL authors="vladimirp"
#
#ENTRYPOINT ["top", "-b"]

FROM golang:1.21.0-alpine3.18 AS builder

# Add GCC compiler in container
RUN apk add build-base

WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -a -o myapp cmd/url-shorter/main.go

# Stage 2: Final stage
FROM alpine:3.18.0

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built binary from the previous stage
COPY --from=builder /app/myapp ./

# Copy config files from previous stage
COPY --from=builder /app/config ./config/

# Set the entry point for the container
CMD ["./myapp"]