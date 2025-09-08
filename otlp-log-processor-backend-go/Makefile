# The build target should stay at the top since we want it to be the default target.
.PHONY: build
build:
	go build -o otlp-log-processor-backend ./main.go

.PHONY: test
test:
	go test $(GOTEST_FLAGS) -race ./...

.PHONY: fmt
fmt:
	gofumpt -l -w .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: run
run:
	go run ./main.go

.PHONY: install-tools
install-tools:
	@echo Installing tools from tools/go.mod
	@go list -modfile=tools/go.mod tool | xargs -I % go list -modfile=tools/go.mod -f "%@{{.Module.Version}}" % | xargs -tI % go install %

.PHONY: generate
generate:
	go generate -x ./...
