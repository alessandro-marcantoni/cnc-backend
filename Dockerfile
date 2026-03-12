# ---------- Build stage ----------
FROM --platform=$BUILDPLATFORM golang:1.25.5-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

ARG TARGETOS
ARG TARGETARCH

# Build static binary
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o app main/main.go


# ---------- Runtime stage ----------
FROM debian:12-slim

WORKDIR /app

# Install wkhtmltopdf dependencies and package
ARG TARGETARCH

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    wget \
    fontconfig \
    libfreetype6 \
    libjpeg62-turbo \
    libpng16-16 \
    libx11-6 \
    libxcb1 \
    libxext6 \
    libxrender1 \
    xfonts-75dpi \
    xfonts-base \
    && rm -rf /var/lib/apt/lists/*

# Install wkhtmltopdf (multi-arch)
RUN set -eux; \
    if [ "$TARGETARCH" = "arm64" ]; then \
        URL="https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-3/wkhtmltox_0.12.6.1-3.bookworm_arm64.deb"; \
    else \
        URL="https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-3/wkhtmltox_0.12.6.1-3.bookworm_amd64.deb"; \
    fi; \
    wget -q $URL -O wkhtmltox.deb; \
    apt-get update; \
    apt-get install -y --no-install-recommends ./wkhtmltox.deb; \
    rm wkhtmltox.deb; \
    rm -rf /var/lib/apt/lists/*


# Copy application
COPY --from=builder /src/app /app/app
COPY --from=builder /src/db /app/db


# Run as non-root
RUN useradd -u 1000 -r appuser
USER appuser


EXPOSE 8080

ENTRYPOINT ["/app/app"]
