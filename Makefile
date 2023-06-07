OPENAPI_FILE=openapi.html
DATABASE_IMAGE=cinema-image
DATABASE_CONTAINER=cinema-container

.PHONY: run
run: clean-docker docker-db
	go run cmd/cinema.go

docker-compose-run:
	docker compose up

.PHONY: test
test:
	go test ./...

clean:
	docker compose down
	docker

.PHONY: docker-db
docker-db:
	docker build -f database/Dockerfile -t $(DATABASE_IMAGE) .
	docker run --name $(DATABASE_CONTAINER) -d -p 5432:5432 $(DATABASE_IMAGE)

.PHONY: clean-docker
clean-docker:
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
