templ:
	@templ generate -watch -proxy=http://localhost:7307

build:
	@templ generate view
	@go build -o ./bin/main ./main.go
