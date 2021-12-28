export TEST_CONTAINER_NAME=postgresql

export JWT_SIGNING_KEY=secret

export DB_NAME=project
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_HOST=0.0.0.0
export DB_PORT=5432
export DB_SSL_MODE=disable

test.integration:
	docker-compose -f docker-compose.test.yml up --build -d
	go test -v ./tests -timeout 30s
	docker stop $$TEST_CONTAINER_NAME