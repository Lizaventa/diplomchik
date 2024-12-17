FROM golang:1.22.10-alpine
WORKDIR /app
# Установка зависимостей
COPY go.mod go.sum ./
RUN go mod download
# Копирование исходного кода
COPY . .
# Сборка
RUN go build -o main .
# Открытие порта
EXPOSE 8080
CMD ["./main"]
