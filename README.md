# ssubench

Платформа для размещения заданий с оплатой виртуальными баллами.

**Архитектура:** Go + PostgreSQL + JWT  
**Роли:** Customer (заказчик), Executor (исполнитель), Admin (администратор)

## Быстрый старт

### 1. Запуск PostgreSQL через Docker

```bash
docker compose up -d
```

Эта команда:
- Запускает PostgreSQL 16 на порту `5432`
- Создаёт БД `ssubench`
- Автоматически выполняет SQL-миграции из `migrations/001_init.sql`

### 2. Установка зависимостей Go

```bash
go mod tidy
```

### 3. Запуск сервера

```bash
go run ./cmd/app/
```

Сервер запустится на `http://localhost:8080`.

## Конфигурация (.env)

Файл `.env` в корне проекта:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ssubench
JWT_SECRET=super_secret_key_2026
HTTP_ADDR=:8080
LOG_LEVEL=INFO
SHUTDOWN_TIMEOUT=10s
```

## Миграции

SQL-миграции находятся в `migrations/001_init.sql`.  
При запуске через Docker Compose файл монтируется в `/docker-entrypoint-initdb.d/` и выполняется автоматически при первом старте контейнера.

Таблицы:
- **users** — пользователи (email, пароль (bcrypt), роль, баланс, блокировка)
- **tasks** — задачи (заказчик, название, описание, бюджет, статус)
- **bids** — отклики (задача, исполнитель, сумма, статус)
- **payments** — платежи (задача, от кого, кому, сумма)



## API Эндпоинты

| Метод | Путь | Роль | Описание |
|-------|------|------|----------|
| POST | `/register` | — | Регистрация пользователя |
| POST | `/login` | — | Получение JWT-токена |
| POST | `/tasks` | Customer | Создание задачи |
| GET | `/tasks` | Любая | Список задач (с пагинацией) |
| POST | `/bids` | Executor | Отклик на задачу |
| POST | `/tasks/{id}/accept-bid` | Customer | Принятие отклика → `in_progress` |
| POST | `/tasks/{id}/mark-completed` | Executor | Задача выполнена → `completed` |
| POST | `/tasks/{id}/confirm` | Customer | Подтверждение + оплата → `confirmed` |

## Примеры curl

### Регистрация пользователя:
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@test.com","password":"pass123","role":"customer"}'
```

### Аутентификация:
```bash
curl -Method POST http://localhost:8080/login `
-Headers @{ "Content-Type" = "application/json" } `
-Body '{"email":"test@test.com","password":"mysecretpassword"}'
```

### Создание задачи заказчиком:
```bash
curl -Method POST http://localhost:8080/tasks `
-Headers @{ 
    "Authorization" = "Bearer CUSTOMER_TOKEN"
    "Content-Type" = "application/json"
} `
-Body '{"title":"Разработать API","description":"Нужен бэкенд на Go","budget":1000}'```
```

### Отклик на задачу исполнителем:
```bash
curl -Method POST http://localhost:8080/bids `
-Headers @{ 
    "Authorization" = "Bearer EXECUTOR_TOKEN"
    "Content-Type" = "application/json"
} `
-Body '{"task_id":1,"amount":900}'
```

### Принять отклик
```bash
curl -Method POST "http://localhost:8080/tasks/1/accept-bid" `
-Headers @{ 
    "Authorization" = "Bearer CUSTOMER_TOKEN"
    "Content-Type" = "application/json"
} `
-Body '{"bid_id":1}'
```

### Отметить как выполненное
```bash
curl -Method POST "http://localhost:8080/tasks/1/mark-completed" `
-Headers @{ 
    "Authorization" = "Bearer EXECUTOR_TOKEN"
}
```

### Подтверждение + оплата
```bash
curl -Method POST "http://localhost:8080/tasks/1/confirm" `
-Headers @{ 
    "Authorization" = "Bearer CUSTOMER_TOKEN"
}
```
