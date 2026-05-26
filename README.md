# Go Backend Template

Минимальный шаблон Go-сервиса:
- cmd/app/main.go
- internal/ (config, httpserver)
- /health
- slog JSON logs
- graceful shutdown
- tests
- Makefile
- Docker

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


## Lecture 2 examples
```bash
# runnable demo
go run ./cmd/examples

# tests
go test ./...
```


## Lecture 2 practice endpoints
- GET /demo/zero
- GET /demo/sum?a=10&b=20
- GET /demo/switch?x=0
- POST /demo/bytes

See PRACTICE.md for details.


## Lecture 3 demo endpoints
- GET /demo3/parse?x=...
- GET /demo3/user?id=...
- GET /demo3/wrap
- GET /demo3/slice
- GET /demo3/map

See PRACTICE3.md for details.
