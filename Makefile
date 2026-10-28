BINARY := borlabd
GOFLAGS ?=

.PHONY: build test lint fmt tidy clean

build:
	go build $(GOFLAGS) -o bin/$(BINARY) ./cmd/borlabd

test:
	go test $(GOFLAGS) ./...

lint:
	go vet ./...

fmt:
	gofmt -l -w .

tidy:
	go mod tidy

clean:
	rm -rf bin
