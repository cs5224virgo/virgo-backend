# Stage 1: Build the application
# Use the official Golang image as a builder base
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
# This is done before copying the source code to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Compile the application to a single binary
RUN CGO_ENABLED=0 GOOS=linux go build -o virgo .

FROM alpine AS base
WORKDIR /app
RUN addgroup -g 1000 virgo && \
    adduser -D -u 1000 -G virgo virgo && \
    apk update && apk add ca-certificates
COPY --chown=1000:1000 --from=builder /app/virgo /app/virgo
COPY ./config.yaml /app
COPY ./db/migrations /app/db/migrations
EXPOSE 7000

ENTRYPOINT ["/app/virgo"]