GO = go
FIRESTORE_EMULATOR_HOST = localhost:8833
GOOGLE_CLOUD_PROJECT = emulator

.PHONY: run
run:
	$(GO) run ./cmd/linebot

.PHONY: run-with-emulator
	FIRESTORE_EMULATOR_HOST="$(FIRESTORE_EMULATOR_HOST)" \
	GOOGLE_CLOUD_PROJECT="$(GOOGLE_CLOUD_PROJECT)" \
	$(GO) run ./cmd/linebot

.PHONY: generate
generate:
	GOFLAGS=-mod=mod $(GO) generate ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: emulator
emulator:
	firebase emulators:start --project="$(GOOGLE_CLOUD_PROJECT)"
