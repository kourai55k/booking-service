.PHONY: start docker-build docker-run help docker-stop docker-rm compose-up compose-down compose-delete

help:
	@echo Usage:
	@echo   make start - Run the application locally
	@echo   make docker-build - Build docker image
	@echo   make docker-run - Run docker container
	@echo   make compose-up - Run docker-compose
	@echo   make compose-down - Stop and remove docker-compose
	@echo   make compose-delete - Remove docker images

run, r:
	go run cmd/booking-service/main.go

run-local, rl:
	go run cmd/local/main.go

docker-build, db:
	docker build -t booking-service .

docker-run, dr:
	docker run -p 8080:8080 --name=booking-service booking-service

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

compose-delete:
	docker rmi booking-service-booking-service:latest
	docker rmi postgres:16