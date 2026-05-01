.PHONY: dev docs tidy

# Run the API
dev:
	go run ./cmd/api

# Generate Swagger docs (pointing to the new main location)
docs:
	swag init -g cmd/api/main.go --parseInternal --parseDependency --dir ./ --output ./docs

# Clean up dependencies
tidy:
	go mod tidy