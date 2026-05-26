APP_NAME=app
PKG=./...

.PHONY: fmt test race bench run build

fmt:
	gofmt -w .

test:
	go test $(PKG)

race:
	go test -race $(PKG)

bench:
	go test -bench=. $(PKG)

run:
	go run ./cmd/app

build:
	go build -o bin/$(APP_NAME) ./cmd/app
