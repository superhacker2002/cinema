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

.PHONY: test-auth
test-auth:
	 curl -i -X 'POST' \
    'localhost:8080/auth/' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{ \
    "username": "test_user", \
    "password": "10a6e6cc8311a3e2bcc09bf6c199adecd5dd59408c343e926b129c4914f3cb01" \
    }'

.PHONY: test-create-user
test-create-user:
	 curl -i -X 'POST' \
    'localhost:8080/users/' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{ \
    "username": "test_user", \
    "password": "10a6e6cc8311a3e2bcc09bf6c199adecd5dd59408c343e926b129c4914f3cb01" \
    }'

.PHONY: test-get-sessions
test-get-sessions:
	 curl -i -X 'GET' \
    'localhost:8080/cinema-sessions/1?date=2023-05-29' \
    -H 'accept: application/json'

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

.PHONY: test-get-halls
test-get-halls:
	 curl -i -X 'GET' \
    'localhost:8080/halls/' \
    -H 'accept: application/json'

test-get-hall:
	 curl -i -X 'GET' \
    'localhost:8080/halls/2' \
    -H 'accept: application/json'

test-create-hall:
	 curl -i -X 'POST' \
    'localhost:8080/halls/' \
    -H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{ \
	"name": "lol", \
	"capacity": 100 \
	}'

test-delete-hall:
	 curl -i -X 'DELETE' \
    'localhost:8080/halls/1'


test-seats:
	 curl -i -X 'GET' \
    'localhost:8080/cinema-sessions/1/seats'