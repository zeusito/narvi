# syntax=docker/dockerfile:1

# -------------------------------------------------------------
# Stage 1: Compiling Application Binary (Builder)
# -------------------------------------------------------------
FROM golang:1.26-alpine AS builder
WORKDIR /src

# Install critical system tools (certificates, timezone database, git)
RUN apk add --no-cache git ca-certificates tzdata

# Copy dependency manifests first to leverage Docker layer caching
COPY go.mod go.sum ./

# Download Go module dependencies using BuildKit cache mounts to accelerate subsequent builds
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the complete codebase
COPY . .

# Statically compile the Go binary for Linux:
# - CGO_ENABLED=0 ensures the binary is fully statically linked (doesn't require glibc)
# - ldflags "-s -w" strips debug symbols and DWARF tables to minimize binary footprint
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/narvi ./cmd/main.go

# -------------------------------------------------------------
# Stage 2: Production Execution Environment
# -------------------------------------------------------------
# distroless/static-debian13:nonroot contains only minimal standard runtime environment,
# offering the absolute minimum attack surface area.
FROM gcr.io/distroless/static-debian13:nonroot AS release

# Establish application directory
WORKDIR /app

# Safely copy CA certificates and timezone database compiled in the builder image
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy statically compiled Go binary from builder stage with correct nonroot ownership
COPY --from=builder --chown=nonroot:nonroot /app/narvi /app/narvi

# Expose the application listening port
EXPOSE 8080

# Enforce running container under unprivileged 'nonroot' UID (65532)
USER nonroot

# Run the compiled API binary
ENTRYPOINT ["/app/narvi"]
