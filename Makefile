OPENAPI_FILE := openapi.html

.PHONY: run
run:
	docker compose up -d

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

