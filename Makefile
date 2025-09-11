.PHONY: up

up:
	docker-compose build --no-cache pocketbase backend && \
	docker-compose up pocketbase backend
