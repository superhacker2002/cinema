OPENAPI_FILE=openapi.html
DATABASE_CONTAINER=cinema-container
DATABASE_IMAGE=cinema-image

.PHONY: run
run: docker_db
	go run cmd/cinema.go

.PHONY: test
test:
	go test ./...

.PHONY: test_auth
test_auth:
	 curl -i -X 'POST' \
    'localhost:8080/auth/login/' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{ \
    "username": "admin", \
    "password": "password" \
    }'

.PHONY: docker_db
docker_db:
	docker build -f database/Dockerfile -t $(DATABASE_IMAGE) .
	docker run --name $(DATABASE_CONTAINER) -d -p 5432:5432 $(DATABASE_IMAGE)

.PHONY: clean_docker
clean_docker:
	docker stop $(DATABASE_CONTAINER)
	docker rm $(DATABASE_CONTAINER)
	docker rmi $(DATABASE_IMAGE)

.PHONY: openapi-docs
openapi-docs:
	redocly build-docs api/openapi.yaml --output=docs/$(OPENAPI_FILE)
	@echo "Open html file created in docs/ directory with the browser."

.PHONY: clean-docs
clean-docs:
	rm -rf docs/


