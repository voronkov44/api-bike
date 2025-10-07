# Stage 1: Build
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates bash

WORKDIR /src

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org && go mod download

# Копируем весь код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/app ./cmd

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/migrate ./migrations

# Stage 2: Run
FROM alpine:3.18
RUN apk add --no-cache ca-certificates bash

# Создаем пользователя
RUN addgroup -S app && adduser -S -G app app

# Копируем бинарник и скрипт запуска
COPY --from=builder /app/bin/app /usr/local/bin/app
COPY --from=builder /app/bin/migrate /usr/local/bin/migrate
COPY start.sh /usr/local/bin/start.sh
RUN chmod +x /usr/local/bin/start.sh
RUN chown app /usr/local/bin/app /usr/local/bin/migrate /usr/local/bin/start.sh

USER app
EXPOSE 8081
ENV PORT=8081

ENTRYPOINT ["/usr/local/bin/start.sh"]
