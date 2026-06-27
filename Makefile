ifeq ($(OS),Windows_NT)
    BIN := artisan.exe
else
    BIN := artisan
endif

.PHONY: build serve serve-dev swagger migrate-up migrate-down

build:
	go build -o $(BIN) artisan.go

serve: build
	./$(BIN) serve

serve-dev: build
	./$(BIN) serve:dev

swagger: build
	./$(BIN) swagger:generate

migrate-up: build
	./$(BIN) migrate:up

migrate-down: build
	./$(BIN) migrate:down 1
