default: fmt build

fmt:
	go fmt github.com/DimensionDataResearch/go-octo-api/...

build:
	go build

test:
	go test -v github.com/DimensionDataResearch/go-octo-api
