.PHONY: all run

run:
	docker-compose up -d --build
	docker image prune -f