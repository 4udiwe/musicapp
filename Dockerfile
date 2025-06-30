# Используем официальный образ Go
FROM golang:1.24 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# 1. Копируем только файлы модулей для кэширования зависимостей
COPY go.mod go.sum ./

# 2. Скачиваем зависимости с кэшированием (используем BuildKit cache mount)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# 3. Копируем остальные файлы проекта
COPY . .

# 4. Опционально: создаем vendor-папку для полного кэширования зависимостей
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod vendor

# 5. Собираем приложение (статически линкуем с использованием vendor)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /app/main cmd/main.go

# Используем ultra-легковесный образ для запуска
FROM alpine:3.19

# Устанавливаем зависимости для alpine (если нужны)
RUN apk add --no-cache ca-certificates tzdata

# Копируем бинарник из builder-этапа
COPY --from=builder /app/main /app/main

# Копируем статические файлы/конфиги если есть (пример)
# COPY --from=builder /app/static /app/static
# COPY --from=builder /app/config /app/config

# Открываем порт, который слушает приложение
EXPOSE 8080

# Указываем рабочую директорию при запуске
WORKDIR /app

# Запускаем приложение
CMD ["/app/main"]