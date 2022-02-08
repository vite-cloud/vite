.PHONY: test checks build

tests: test
test:
	go test -v -count=1 ./...
checks:
	gosec ./...
	go fmt ./...
	go vet ./...
	golangci-lint run	 ./...
build:
	go build -ldflags='-w -s -X github.com/redwebcreation/nest/globals.Version=testing' -gcflags=all='-l'