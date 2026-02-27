# Build stage
FROM debian:trixie AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    golang-go \
    pkg-config \
    libopus-dev \
    libopusfile-dev \
    git \
    ca-certificates \
    gcc \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

# Copy go.mod and go.sum files first for better layer caching
COPY webrtc-server/go.mod webrtc-server/go.sum ./webrtc-server/
COPY webrtc-agents/go.mod webrtc-agents/go.sum ./webrtc-agents/

# Download dependencies
WORKDIR /build

# Copy source code
COPY webrtc-server/ ./webrtc-server/
COPY webrtc-agents/ ./webrtc-agents/

WORKDIR /build/webrtc-agents
RUN go mod download

# Build the binary
ARG TARGETARCH
RUN \
    mkdir -p /build/webrtc-agents/bin  && \
    CGO_ENABLED=1 GOOS=linux GOARCH=$TARGETARCH go build -ldflags="-s -w" -o /build/webrtc-agents/bin/agents ./cmd/agents

# Runtime stage
FROM debian:trixie-slim

# Install runtime dependencies (only what's needed to run the binary)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libopus0 \
    libopusfile0 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --create-home --shell /bin/bash appuser

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/webrtc-agents/bin/agents /app/agents

# Use non-root user for security
USER appuser

# Default environment variables (can be overridden at runtime)
ENV WS_SERVER=ws://localhost:3001/ws

ENTRYPOINT ["/app/agents"]
