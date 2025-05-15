up:
	docker compose up

down:
	docker compose down

clean:
	docker compose down --volumes

start-producer:
	./build/producer
start-uploader:
	./build/uploader
start-presenter:
	./build/presenter

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
	rm -rf ./build/producer
	go build -o ./build/producer cmd/producer/producer.go
build-uploader:
	rm -rf ./build/uploader
	go build -o ./build/uploader cmd/uploader/uploader.go
build-presenter:
	rm -rf ./build/presenter
	go build -o ./build/presenter cmd/presenter/presenter.go

fmt:
	go mod tidy
	go fmt ./...
