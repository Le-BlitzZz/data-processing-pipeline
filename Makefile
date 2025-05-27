up:
	docker compose up
down:
	docker compose down
up-fresh:
	docker compose up --build
clean:
	docker compose down --volumes

lock-processor:
	cd services/processor && poetry lock

start-publisher:
	./bin/publisher
start-dataserver:
	./bin/dataserver
start-processor:
	cd services/processor && poetry run processor
start-trainer:
	cd services/trainer && poetry run trainer


terminal-publisher:
	docker compose exec -u root publisher bash
terminal-dataserver:
	docker compose exec -u root dataserver bash
terminal-processor:
	docker compose exec -u root processor bash
terminal-trainer:
	docker compose exec -u root trainer bash

build: build-publisher build-dataserver
build-publisher:
	rm -rf ./bin/publisher
	go build -o ./bin/publisher cmd/publisher/publisher.go
build-dataserver:
	rm -rf ./bin/dataserver
	go build -o ./bin/dataserver cmd/dataserver/dataserver.go
build-processor:
	cd services/processor && poetry install
build-trainer:
	cd services/trainer && poetry install

fmt:
	go mod tidy
	go fmt ./...
