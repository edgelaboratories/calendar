COVERAGE_OUT ?= bin/cover.out

.PHONY: fmt
fmt:
	go mod tidy
	gofumpt -w .

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	mkdir -p bin
	go test -race -tags=integration -coverprofile=$(COVERAGE_OUT) ./...
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_OUT:.out=.html)

.PHONY: benchmark
benchmark:
	go test -tags=integration -benchmem ./... -run=Bench -bench=.
