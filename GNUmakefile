default: testacc

# Build and install the provider to GOBIN
.PHONY: install
install:
	go install .

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Run CI emulation locally
.PHONY: test-ci
test-ci:
	./scripts/test-ci.sh
