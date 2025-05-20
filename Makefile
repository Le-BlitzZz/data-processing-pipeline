up:
	docker compose up

down:
	docker compose down

up-fresh:
	docker compose up --build

clean:
	docker compose down --volumes

start-producer:
	./bin/producer
start-uploader:
	./bin/uploader
start-presenter:
	./bin/presenter

terminal-producer:
	docker compose exec -u root producer bash
terminal-processor:
	docker compose exec -u root processor bash
terminal-uploader:
	docker compose exec -u root uploader bash
terminal-presenter:
	docker compose exec -u root presenter bash

build: build-producer build-uploader build-presenter
build-producer:
	rm -rf ./bin/producer
	go build -o ./bin/producer cmd/producer/producer.go
build-uploader:
	rm -rf ./bin/uploader
	go build -o ./bin/uploader cmd/uploader/uploader.go
build-presenter:
	rm -rf ./bin/presenter
	go build -o ./bin/presenter cmd/presenter/presenter.go

fmt:
	go mod tidy
	go fmt ./...
