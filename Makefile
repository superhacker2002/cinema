OPENAPI_FILE=openapi.html

.PHONY: run
<<<<<<< HEAD
run: docker-db
=======
run:
>>>>>>> 8426025b7a7b0bd59d3fee9c5bd7960d4b6c6c52
	go run cmd/cinema.go

docker-compose-run:
	docker compose up

.PHONY: test
test:
	go test ./...

<<<<<<< HEAD
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
    'localhost:8080/cinema-sessions/2?date=2023-05-18' \
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

=======
>>>>>>> 8426025b7a7b0bd59d3fee9c5bd7960d4b6c6c52
.PHONY: openapi-docs
openapi-docs:
	redocly build-docs api/openapi.yaml --output=docs/$(OPENAPI_FILE)
	@echo "Open html file created in docs/ directory with the browser."

.PHONY: clean-docs
clean-docs:
<<<<<<< HEAD
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
=======
	rm -rf docs/
>>>>>>> 8426025b7a7b0bd59d3fee9c5bd7960d4b6c6c52
