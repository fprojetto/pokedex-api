# ------------------------------------------------------------------------------
# STAGE 1: Build the Binary
# ------------------------------------------------------------------------------
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# 1. Cache Dependencies
# Copy only go.mod and go.sum first. Docker will cache this layer
# so you don't re-download dependencies if only your code changes.
COPY go.mod ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build flags explanation:
# CGO_ENABLED=0 : Disables CGO to create a statically linked binary (no external C library deps)
# GOOS=linux    : Ensures we build for Linux (even if you build on Mac/Windows)
# -ldflags="-w -s" : Strips DWARF debug info and symbol tables to reduce binary size by ~30%
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/api

# ------------------------------------------------------------------------------
# STAGE 2: The Production Image
# ------------------------------------------------------------------------------
# We use Alpine for a balance of small size (~5MB) and debuggability (it has a shell)
FROM alpine:latest AS runner

# Install CA certificates (essential for making HTTPS requests to other APIs)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/server .

# Create a non-root user for security
# It is bad practice to run as root in production
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Expose the port your app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./server"]
