
build:
	rm -rf out && go build -o out/ src/main.go

dev:
	go run src/main.go

test:
	go test -coverprofile=coverage.out ./...
