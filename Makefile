build:
	@go build -o main ./cmd/server

run: build
	@./main

docker-run:
	@docker compose up

docker-run-build:
	@docker compose up --build

docker-stop:
	@docker compose down

docker-stop-v:
	@docker compose down -v