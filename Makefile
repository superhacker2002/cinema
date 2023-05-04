OPENAPI_FILE := openapi.html

.PHONY: run
run:
	go run cmd/cinema.go

.PHONY: test
test:
	go test ./...

.PHONY: test_auth
test_auth:
	 curl -i -X 'POST' \
    'localhost:8080/auth/login' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{ \
    "username": "invalid", \
    "password": "password" \
    }'

.PHONY: docker_db
docker_db:
	docker build -f database/Dockerfile -t cinema-image .
	docker run --name cinema-container -d -p 5432:5432 cinema-image

.PHONY: clean_docker
clean_docker:
	docker stop cinema-container
	docker rm cinema-container
	docker rmi cinema-image



.PHONY: openapi-docs
openapi-docs:
	redocly build-docs api/openapi.yaml --output=docs/$(OPENAPI_FILE)
	@echo "Open html file created in docs/ directory with the browser."

.PHONY: clean-docs
clean-docs:
	rm -rf docs/


