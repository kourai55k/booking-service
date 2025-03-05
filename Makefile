.PHONY: start docker-build docker-run help docker-stop docker-rm

help:
	@echo Usage:
	@echo   make start - Run the application locally
	@echo   make docker-build - Build docker image
	@echo   make docker-run - Run docker container

run, r:
	go run cmd/booking-service/main.go

docker-build, db:
	docker build -t booking-service .

docker-run, dr:
	docker run -p 8080:8080 --name=booking-service booking-service
