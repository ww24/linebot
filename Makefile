GO = go
FIRESTORE_EMULATOR_HOST = :8833
GOOGLE_CLOUD_PROJECT = emulator
TESTTIME = go run github.com/tenntenn/testtime/cmd/testtime

.PHONY: run
run:
	$(GO) run ./cmd/linebot

.PHONY: run-with-emulator
	FIRESTORE_EMULATOR_HOST="$(FIRESTORE_EMULATOR_HOST)" \
	GOOGLE_CLOUD_PROJECT="$(GOOGLE_CLOUD_PROJECT)" \
	$(GO) run ./cmd/linebot



.PHONY: generate
generate:
	$(GO) generate ./...

.PHONY: test
test:
	$(GO) test -v --race -overlay=`$(TESTTIME)` ./...

.PHONY: integration-test
integration-test:
	FIRESTORE_EMULATOR_HOST="$(FIRESTORE_EMULATOR_HOST)" \
	$(GO) test -v --race -overlay=`$(TESTTIME)` -tags=integration ./infra/...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: emulator
emulator:
	firebase emulators:start --project="$(GOOGLE_CLOUD_PROJECT)"
