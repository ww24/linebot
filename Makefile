GO = go

.PHONY: run
run:
	$(GO) run ./cmd/linebot

.PHONY: generate
generate:
	$(GO) generate ./...

.PHONY: lint
lint:
	golangci-lint run
