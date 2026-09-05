.PHONY: templ build build-exampleserver test run-example

# Regenerate templ components and hot-reload them against the running web UI.
templ:
	@templ generate -watch -proxy=http://localhost:7307

# Build the main CLI into ./bin/hextest
build:
	@templ generate
	@go build -o ./bin/hextest ./cmd/hextest

# Build the standalone demo server into ./bin/exampleserver
build-exampleserver:
	@go build -o ./bin/exampleserver ./cmd/exampleserver

test:
	@go test ./...

run-example:
	@go run ./cmd/exampleserver
