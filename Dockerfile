FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Копируем только go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org && go mod download

# Копируем весь код
COPY . .

# Сборка бинарника приложения
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/app ./cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates

# Создаем пользователя (не root)
RUN addgroup -S app && adduser -S -G app app

# Копируем бинарник
COPY --from=builder /app/bin/app /usr/local/bin/app
RUN chown app /usr/local/bin/app
USER app

EXPOSE 8081
ENV PORT=8081

ENTRYPOINT ["/usr/local/bin/app"]
