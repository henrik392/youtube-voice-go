# Stage 1: Build the application
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o main cmd/api/main.go

# Stage 2: Run the application
FROM alpine:latest

WORKDIR /app

# Install yt-dlp
RUN apk update && apk add --no-cache yt-dlp ffmpeg ca-certificates

RUN mkdir -p /app/downloads
# RUN chmod 777 /app/downloads

COPY --from=builder /app/main .
# COPY --from=builder /app/.env .
# COPY static/ ./static/           # Copy static files if you have any

EXPOSE 8080

CMD ["./main"]