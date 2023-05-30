OPENAPI_FILE=openapi.html
DATABASE_CONTAINER=cinema-container
DATABASE_IMAGE=cinema-image

.PHONY: run
run: clean-docker docker-db
	go run cmd/cinema.go

.PHONY: test
test:
	go test ./...

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
	 curl -s -i -X 'POST' \
    'localhost:8080/users/' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{ \
    "username": "test_user, \
    "password": "10a6e6cc8311a3e2bcc09bf6c199adecd5dd59408c343e926b129c4914f3cb01" \
    }'

.PHONY: test-get-sessions
test-get-sessions:
	curl -X 'GET' \
        'localhost:8080/cinema-sessions/2?date=2023-05-22' \
        -H 'accept: application/json' | jq .

.PHONY: test-get-all-sessions
test-get-all-sessions:
	curl -X 'GET' \
        'localhost:8080/cinema-sessions/' \
        -H 'accept: application/json' | jq .


.PHONY: test-delete-session
test-delete-session:
	curl -X 'DELETE' \
        'localhost:8080/cinema-sessions/5'

.PHONY: test-create-session
test-create-session:
	curl -i -X 'POST' \
        'localhost:8080/cinema-sessions/2' \
		-H 'accept: application/json' \
		-H 'Content-Type: application/json' \
		-d '{ \
		"movieId": 1, \
		"startTime": "2023-05-30 17:30:00 +04", \
		"price": 10.0 \
		}'

.PHONY: test-update-session
test-update-session:
	curl -i -X 'PUT' \
        'localhost:8080/cinema-sessions/2' \
		-H 'accept: application/json' \
		-H 'Content-Type: application/json' \
		-d '{ \
		"movieId": 1, \
		"startTime": "2023-05-30 17:30:00 +04", \
		"price": 20.0 \
		}'

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


