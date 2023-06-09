OPENAPI_FILE=openapi.html
DATABASE_CONTAINER=cinema-container
DATABASE_IMAGE=cinema-image
MINIO_ACCESS_KEY=rubiezzy
MINIO_SECRET_KEY=a@3JsY4Vn&fT8s
MINIO_CONTAINER=minio

.PHONY: run
run: clean-db docker-db
	go run cmd/cinema.go

.PHONY: test
test:
	go test ./...

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
      -e "MINIO_ACCESS_KEY=$(MINIO_ACCESS_KEY)" \
      -e "MINIO_SECRET_KEY=$(MINIO_SECRET_KEY)" \
      minio/minio server /data \
      mc mb $(MINIO_CONTAINER)/tickets


.PHONY: openapi-docs
openapi-docs:
	redocly build-docs api/openapi.yaml --output=docs/$(OPENAPI_FILE)
	@echo "Open html file created in docs/ directory with the browser."

.PHONY: clean-docs
clean-docs:
	rm -rf docs/

