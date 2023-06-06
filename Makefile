OPENAPI_FILE=openapi.html

.PHONY: run
run:
	docker compose up

.PHONY: openapi-docs
openapi-docs:
	redocly build-docs api/openapi.yaml --output=docs/$(OPENAPI_FILE)
	@echo "Open html file created in docs/ directory with the browser."

.PHONY: clean
clean:
	docker compose down

.PHONY: clean-docs
clean-docs:
	rm -rf docs/

clean-images:
	docker rmi cinemago-app-cinema-service


.PHONY: test-auth
test-auth:
	 curl -i -X 'POST' \
    'localhost:8080/auth/' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -d '{ \
    "username": "admin", \
    "password": "6D4525C2A21F9BE1CCA9E41F3AA402E0765EE5FCC3E7FEA34A169B1730AE386E" \
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
	 curl -i -X 'GET' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODYxMjIzMzMsInVzZXJfaWQiOjF9.fVF1xUBEBNnUNZVU-FbL5aQGYUDrj3QlUcvRoNYTi4Q" \
    'localhost:8080/cinema-sessions/1?date=2023-05-29' \
    -H 'accept: application/json'

.PHONY: test-create-session
test-create-session:
	 curl curl -i -X 'POST' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODYxMjIzMzMsInVzZXJfaWQiOjF9.fVF1xUBEBNnUNZVU-FbL5aQGYUDrj3QlUcvRoNYTi4Q" \
    'localhost:8080/cinema-sessions/2' \
    -H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{ \
	"movieId": 3, \
	"startTime": "2023-05-29 14:00:00 +04", \
	"price": 10.0 \
	}'

.PHONY: test-update-session
test-update-session:
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

test-delete-session:
	 curl -i -X 'DELETE' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODYxMjIxODcsInVzZXJfaWQiOjJ9.5BfTo-05IIPJVlAqpFDXplfRCNzkkneOIuu6nIjnXnE" \
 	'localhost:8080/cinema-sessions/100000'

.PHONY: test-get-halls
test-get-halls:
	 curl -i -X 'GET' \
    'localhost:8080/halls/' \
    -H 'accept: application/json'

test-get-hall:
	 curl -i -X 'GET' \
    'localhost:8080/halls/1' \
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

test-create-ticket:
	 curl -i -X 'POST' \
 	-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODYxMjgyMDYsInVzZXJfaWQiOjF9.wD8AOseT-5wYboBkQwL_BChShjCJN7nB7cze6A8izyI" \
    'localhost:8080/tickets/'

