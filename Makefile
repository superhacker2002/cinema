OPENAPI_FILE=openapi.html
DATABASE_IMAGE=cinema-image
DATABASE_CONTAINER=cinema-container
MINIO_ROOT_USER=rubiezzy
MINIO_ROOT_PASSWORD=a3JsY4VnfT8s
MINIO_CONTAINER=minio

.PHONY: run
run:
	docker compose --env-file .env up

.PHONY: test
test:
	go test ./...

clean:
	docker compose down

.PHONY: docker-db
docker-db:
	docker build -f database/Dockerfile -t $(DATABASE_IMAGE) .
	docker run --name $(DATABASE_CONTAINER) -d -p 5432:5432 $(DATABASE_IMAGE)

.PHONY: clean-db
clean-db:
	docker stop $(DATABASE_CONTAINER)
	docker rm $(DATABASE_CONTAINER)
	docker rmi $(DATABASE_IMAGE)

.PHONY: clean-minio
clean-minio:
	docker stop $(MINIO_CONTAINER)
	docker rm $(MINIO_CONTAINER)

.PHONY: minio-storage
minio-storage:
	docker run -p 9000:9000 --name $(MINIO_CONTAINER) \
      -e "MINIO_ROOT_USER=$(MINIO_ROOT_USER)" \
      -e "MINIO_ROOT_PASSWORD=$(MINIO_ROOT_PASSWORD)" \
      minio/minio server /data

.PHONY: create-bucket
create-bucket:
	mc alias set minio http://localhost:9000 $(MINIO_ROOT_USER) $(MINIO_ROOT_PASSWORD)
	mc mb minio/tickets
	mc anonymous set public minio/tickets

.PHONY: openapi-docs
openapi-docs:
	redocly build-docs api/openapi.yaml --output=docs/$(OPENAPI_FILE)
	@echo "Open html file created in docs/ directory with the browser."

.PHONY: clean-docs
clean-docs:
	rm -rf docs/
