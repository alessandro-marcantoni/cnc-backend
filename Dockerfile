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
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/app /app/app

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/app"]
