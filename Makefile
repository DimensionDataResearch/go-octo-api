default: build test

fmt:
	go fmt github.com/DimensionDataResearch/go-octo-api/...

build: fmt
	go build

test: fmt
	go test -v github.com/DimensionDataResearch/go-octo-api
