# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем файлы модулей и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o tesla-app ./cmd/app

# Final stage
FROM alpine:latest

# Устанавливаем необходимые пакеты для работы
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем бинарник из стадии сборки
COPY --from=builder /app/tesla-app .

# Копируем шаблоны и статические файлы
#COPY --from=builder /app/templates ./templates
#COPY --from=builder /app/static ./static

# Копируем конфигурационные файлы из корня проекта
COPY --from=builder /app/config/config.toml .
COPY --from=builder /app/.env .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./tesla-app"]