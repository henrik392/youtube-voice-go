# Stage 1: Build the application
FROM golang:alpine AS builder

# Install Node.js and npm
RUN apk add --no-cache nodejs npm

WORKDIR /app

# Copy Go module files first for better caching
COPY go.mod ./
COPY go.sum ./

RUN go mod download

# Copy package files for Node.js dependencies
COPY package.json ./
COPY package-lock.json ./

# Install Node.js dependencies
RUN npm ci --only=production

COPY . .

# Build CSS with Tailwind
RUN npx tailwindcss -i cmd/web/assets/css/input.css -o cmd/web/assets/css/output.css

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
COPY --from=builder /app/cmd/web/assets ./cmd/web/assets

EXPOSE 8080

CMD ["./main"]