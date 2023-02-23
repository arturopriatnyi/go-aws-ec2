all: build run

build:
	docker-compose build

run:
	docker-compose up

test:
	go test -cover -coverprofile=coverage.html -timeout 30s ./...

.PHONY: coverage
coverage:
	go tool cover -html=coverage.html
