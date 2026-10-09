.PHONY: build vet fmt tidy

build:
	go build -o build/ggb .

vet:
	go vet ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

tidy:
	go mod tidy
