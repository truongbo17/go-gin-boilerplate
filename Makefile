.PHONY: build test vet fmt tidy

build:
	go build -o build/ggb .

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

tidy:
	go mod tidy
