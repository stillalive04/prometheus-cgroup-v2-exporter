# Multi-stage Dockerfile for Prometheus cgroup v2 Exporter

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /src

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o prometheus-cgroup-v2-exporter \
    ./cmd/prometheus-cgroup-v2-exporter

# Runtime stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user
RUN addgroup -g 65534 exporter && \
    adduser -D -u 65534 -G exporter exporter

# Set working directory
WORKDIR /

# Copy binary from builder stage
COPY --from=builder /src/prometheus-cgroup-v2-exporter /usr/local/bin/prometheus-cgroup-v2-exporter

# Create directories for cgroup access
RUN mkdir -p /sys/fs/cgroup /host/proc

# Change ownership
RUN chown -R exporter:exporter /usr/local/bin/prometheus-cgroup-v2-exporter

# Switch to non-root user
USER exporter

# Expose metrics port
EXPOSE 9753

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9753/health || exit 1

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/prometheus-cgroup-v2-exporter"]

# Default command
CMD ["--web.listen-address=:9753"]

# Labels
LABEL maintainer="DevOps Team <devops@example.com>"
LABEL org.opencontainers.image.title="Prometheus cgroup v2 Exporter"
LABEL org.opencontainers.image.description="High-performance Prometheus exporter for cgroup v2 metrics"
LABEL org.opencontainers.image.vendor="Open Source Community"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/stillalive04/prometheus-cgroup-v2-exporter"
LABEL org.opencontainers.image.documentation="https://github.com/stillalive04/prometheus-cgroup-v2-exporter/blob/main/README.md"
