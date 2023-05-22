OPENAPI_FILE=openapi.html

.PHONY: run
run:
	go run cmd/cinema.go

.PHONY: test
test:
	go test ./...

.PHONY: openapi-docs
openapi-docs:
	redocly build-docs api/openapi.yaml --output=docs/$(OPENAPI_FILE)
	@echo "Open html file created in docs/ directory with the browser."

.PHONY: clean-docs
clean-docs:
	rm -rf docs/


