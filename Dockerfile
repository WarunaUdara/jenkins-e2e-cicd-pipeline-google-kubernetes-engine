# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies)
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001

# Create a minimal passwd file for the application
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0: Build static binary (no C dependencies)
# GOOS=linux: Target Linux OS
# -ldflags="-w -s": Strip debug information to reduce binary size
# -a: Force rebuilding of packages
# -installsuffix cgo: Add suffix to package directory
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=1.0.0" \
    -a -installsuffix cgo \
    -o todo-app .

# Stage 2: Runtime
FROM alpine:latest

# Import the user and group files from the builder
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Import the compiled executable from the builder
COPY --from=builder /build/todo-app /app/todo-app

# Copy static files and templates
COPY --from=builder /build/static /app/static
COPY --from=builder /build/templates /app/templates

# Import ca-certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Use the non-root user
USER appuser:appuser

WORKDIR /app

# Expose application port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./todo-app"]
