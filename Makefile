
build:
	rm -rf out && go build -o out/ src/main.go

dev:
	go run src/main.go

repl:
	go run src/main.go --repl

exe:
	go run src/main.go --exe=$(FILE)

test:
	go test -coverprofile=coverage.out ./...
