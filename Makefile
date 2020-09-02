test:
	go test -cover -race ./...

up:
	docker-compose up -d
