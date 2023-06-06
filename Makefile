OPENAPI_FILE=openapi.html
DATABASE_CONTAINER=cinema-container
DATABASE_IMAGE=cinema-image

.PHONY: run
run: clean-docker docker-db
	go run cmd/cinema.go

.PHONY: test
test:
	go test ./...

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

.PHONY: auth
auth:
	 curl -i -X 'POST' \
    'localhost:8080/auth/' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{ \
    "username": "admin", \
    "password": "6D4525C2A21F9BE1CCA9E41F3AA402E0765EE5FCC3E7FEA34A169B1730AE386E" \
    }'

.PHONY: create-user
create-user:
	 curl -i -X 'POST' \
    'localhost:8080/users/' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{ \
    "username": "test_user", \
    "password": "10a6e6cc8311a3e2bcc09bf6c199adecd5dd59408c343e926b129c4914f3cb01" \
    }'

# ---------------- sessions --------------------
.PHONY: get-sessions
get-sessions:
	 curl -i -X 'GET' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODYxMjIzMzMsInVzZXJfaWQiOjF9.fVF1xUBEBNnUNZVU-FbL5aQGYUDrj3QlUcvRoNYTi4Q" \
    'localhost:8080/cinema-sessions/1?date=2023-05-29' \
    -H 'accept: application/json'

.PHONY: create-session
create-session:
	 curl curl -i -X 'POST' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODYxMjIzMzMsInVzZXJfaWQiOjF9.fVF1xUBEBNnUNZVU-FbL5aQGYUDrj3QlUcvRoNYTi4Q" \
    'localhost:8080/cinema-sessions/2' \
    -H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{ \
	"movieId": 3, \
	"startTime": "2023-05-29 14:00:00 +04", \
	"price": 10.0 \
	}'

.PHONY: update-session
update-session:
	 curl -i -X 'PUT' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODYxMjIxODcsInVzZXJfaWQiOjJ9.5BfTo-05IIPJVlAqpFDXplfRCNzkkneOIuu6nIjnXnE" \
    'localhost:8080/cinema-sessions/1' \
    -H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{ \
	"movieId": 3, \
	"hallId": 1, \
	"startTime": "2023-05-29 15:00:00 +04", \
	"price": 10.0 \
	}'

.PHONY: delete-session
delete-session:
	 curl -i -X 'DELETE' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODYxMjIxODcsInVzZXJfaWQiOjJ9.5BfTo-05IIPJVlAqpFDXplfRCNzkkneOIuu6nIjnXnE" \
 	'localhost:8080/cinema-sessions/100000'

.PHONY: available-seats
available-seats:
	 curl -i -X 'GET' \
    'localhost:8080/cinema-sessions/1/seats'

# ------------------- halls --------------------
.PHONY: get-halls
get-halls:
	 curl -i -X 'GET' \
    'localhost:8080/halls/' \
    -H 'accept: application/json'

.PHONY: get-hall
get-hall:
	 curl -i -X 'GET' \
    'localhost:8080/halls/1' \
    -H 'accept: application/json'

.PHONY: create-hall
create-hall:
	 curl -i -X 'POST' \
    'localhost:8080/halls/' \
    -H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{ \
	"name": "lol", \
	"capacity": 100 \
	}'

.PHONY: delete-hall
delete-hall:
	 curl -i -X 'DELETE' \
    'localhost:8080/halls/1'

# ------------------ tickets -------------------
.PHONY: create-ticket
create-ticket:
	 curl -i -X 'POST' \
 	-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODYxMjkzODYsInVzZXJfaWQiOjF9.VbUCJvOL5Oepk24kGIjVteGKljV-WX_4q-Yhcm4i_gY" \
    'localhost:8080/tickets/' \
	-H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{ \
	"sessionId": 1, \
	"seatNumber": 4 \
	}'
