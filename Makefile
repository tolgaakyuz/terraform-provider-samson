TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=samson

all: build

build-test:
	docker-compose build --no-cache terraform-provider-samson-test

test: build-test
	DOCKER_REGISTRY=none docker-compose run \
					terraform-provider-samson-test \
					bash -c './tools/test.sh'

test-i: build-test
	DOCKER_REGISTRY=none docker-compose run \
		terraform-provider-samson-test \
		bash


default: build

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 TF_LOG=debug SAMSON_TOKEN=123 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

.PHONY: build test test-i testacc vet fmt fmtcheck errcheck vendor-status test-compile website website-test
