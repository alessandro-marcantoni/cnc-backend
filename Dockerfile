# ---------- Build stage ----------
FROM --platform=$BUILDPLATFORM golang:1.25.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o app main/main.go


# ---------- Runtime stage ----------
FROM debian:12-slim

WORKDIR /app

# Install wkhtmltopdf and fonts
RUN apt-get update && apt-get install -y \
    wkhtmltopdf \
    fonts-liberation \
    fonts-dejavu-core \
    fonts-noto \
    fonts-noto-cjk \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy application
COPY --from=builder /app/app /app/app
COPY --from=builder /app/db /app/db

# Create non-root user
RUN useradd -r -u 1000 -g root appuser

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/app"]
