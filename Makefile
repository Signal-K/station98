.PHONY: up

up:
	docker-compose build --no-cache pocketbase && docker-compose build --no-cache backend && docker-compose up
