GO = go
BIN := $(abspath ./bin)
FIRESTORE_EMULATOR_HOST ?= localhost:8833
GOOGLE_CLOUD_PROJECT = emulator
GO_ENV ?= GOBIN=$(BIN)

$(BIN)/testtime:
	$(GO_ENV) $(GO) install github.com/tenntenn/testtime/cmd/testtime@v0.2.2

$(BIN)/mockgen:
	$(GO_ENV) $(GO) install github.com/golang/mock/mockgen@v1.6.0

.PHONY: clean-mock
clean-mock:
	$(RM) -r ./mock

.PHONY: run
run:
	$(GO) run ./cmd/linebot

.PHONY: run-with-emulator
run-with-emulator:
	FIRESTORE_EMULATOR_HOST="$(FIRESTORE_EMULATOR_HOST)" \
	GOOGLE_CLOUD_PROJECT="$(GOOGLE_CLOUD_PROJECT)" \
	$(GO) run ./cmd/linebot

.PHONY: generate
generate: $(BIN)/mockgen
generate: clean-mock
	@$(GO_ENV) PATH="${PATH}:$(BIN)" $(GO) generate ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test: FLAGS ?=
test: $(BIN)/testtime
	FIRESTORE_EMULATOR_HOST="$(FIRESTORE_EMULATOR_HOST)" \
	$(GO_ENV) $(GO) test $(FLAGS) -race -overlay="$(shell $(BIN)/testtime -u)" ./...

.PHONY: emulator
emulator:
	firebase emulators:start --project="$(GOOGLE_CLOUD_PROJECT)"
