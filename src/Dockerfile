FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy all source code
COPY . .

# ✅ Build Linux amd64 binary (for ECS)
RUN GOOS=linux GOARCH=amd64 go build -o main .

# Second, minimal runtime image
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/main .

# Ensure binary is executable
RUN chmod +x main

EXPOSE 8080
CMD ["./main"]

