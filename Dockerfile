# -------- Build stage --------
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app main/main.go

# -------- Runtime stage --------
FROM debian:12-slim

WORKDIR /app

# Install Chromium for PDF generation
RUN apt-get update && \
    apt-get install -y \
    chromium \
    chromium-sandbox \
    fonts-liberation \
    libnss3 \
    libatk-bridge2.0-0 \
    libdrm2 \
    libxkbcommon0 \
    libgbm1 \
    libasound2 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set Chrome path environment variable
ENV CHROME_BIN=/usr/bin/chromium

COPY --from=builder /app/app /app/app

EXPOSE 8080

# Create non-root user
RUN useradd -r -u 1000 -g root appuser

USER appuser

ENTRYPOINT ["/app/app"]
