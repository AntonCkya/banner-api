test:
	docker-compose up -d && docker-compose exec app go test ./cmd

run:
	docker-compose up -d
