cover:
	go test -count=1 -coverprofile=coverage.out $(shell go list ./internal/... | grep -E "internal/api|internal/service" | grep -v -E "mocks|middleware|decorator")
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

.PHONY: integration-test

integration-test:
	@echo "Starting integration tests..."
	docker compose -f docker-compose.test.yaml up \
		--abort-on-container-exit \
		--exit-code-from integration \
		--build \
		--remove-orphans
