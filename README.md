# Order Service gRPC Server

Простой gRPC сервис для управления заказами на Go.

![Мем](https://cdn-images-1.readmedium.com/v2/resize:fit:800/1*zInbi1YU_Oemo--ZxgI9Pg.jpeg)

## Быстрый старт

### 1. Установка зависимостей
go mod download

### 2. Настройка конфигурации
Создай файл .env в папке config/:
GRPC_PORT=50051
LOG_LEVEL=info

### 3. Запуск сервера
go run cmd/server/main.go

## gRPC API

### Методы:
- CreateOrder - создание заказа
- GetOrder - получение заказа по ID  
- UpdateOrder - обновление заказа
- DeleteOrder - удаление заказа
- ListOrders - список всех заказов

### Примеры запросов:
# Создать заказ
grpcurl -plaintext -d '{"item": "Laptop", "quantity": 2}' localhost:50051 api.OrderService/CreateOrder

# Получить все заказы
grpcurl -plaintext -d '{}' localhost:50051 api.OrderService/ListOrders

## Конфигурация

Переменная: GRPC_PORT - Порт gRPC сервера - По умолчанию: 50051
Переменная: LOG_LEVEL - Уровень логирования - По умолчанию: info

## Структура проекта

rpc/
- cmd/server/ - Точка входа
- internal/
  - config/ - Конфигурация
  - server/ - Бизнес-логика
  - interceptor/ - gRPC интерсепторы
- pkg/api/test/ - Сгенерированный gRPC код
- config/
  - .env - Конфигурация (не в git)
  - env.example - Пример конфигурации

## Разработка

Сервер автоматически логирует все входящие запросы и ответы через gRPC интерсептор.

Для отладки установи LOG_LEVEL=debug в .env файле.