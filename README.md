# ssubench

## Run
```bash
make run
# GET http://localhost:8080/health
```

## Test
```bash
make test
make race
make bench
```

## Build
```bash
make build
./bin/app
```

## Docker
```bash
docker compose up --build
```


## Регистрация пользователя:
```bash
curl -X POST http://localhost:8080/register
-H "Content-Type: application/json"
-d '{"email": "test@test.com", "password": "mysecretpassword", "role": "customer"}'
```

## Вход:
```bash
curl -X POST http://localhost:8080/login
-H "Content-Type: application/json"
-d '{"email": "test@test.com", "password": "mysecretpassword"}'
```

## Создание задачи заказчиком:
```bash
curl -X POST http://localhost:8080/tasks
-H "Authorization: Bearer CUSTOMER_TOKEN"
-H "Content-Type: application/json"
-d '{"title": "Разработать API", "description": "Нужен бэкенд на Go", "budget": 1000}'
```

## Отклик на задачу исполнителем:
```bash
curl -X POST http://localhost:8080/bids
-H "Authorization: Bearer EXECUTOR_TOKEN"
-H "Content-Type: application/json"
-d '{"task_id": 1, "amount": 900}'
```
