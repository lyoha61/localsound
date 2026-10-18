# Этап сборки бинарника
FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /build/server ./cmd/api

# Финальный легковесный образ
FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache \
    curl \
    ffmpeg \
    nodejs \
    python3 \
    ca-certificates

RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp \
    && chmod a+rx /usr/local/bin/yt-dlp

COPY --from=builder /build/server /app/server

EXPOSE 8080

ENTRYPOINT ["/app/server"]
