# Используем официальный образ Go (легковесный alpine-вариант)
FROM golang:1.24 AS builder

# Копируем исходный код в контейнер
WORKDIR /app
COPY . .

# Скачиваем зависимости (Gin и другие)
RUN go mod download

# Собираем приложение (статически линкуем для уменьшения размера)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Используем ultra-легковесный образ для запуска
FROM alpine

# Копируем бинарник из builder-этапа
COPY --from=builder /app/main /app/main

# Открываем порт, который слушает Gin (по умолчанию 8080)
EXPOSE 8080

# Запускаем приложение
CMD ["/app/main"]
