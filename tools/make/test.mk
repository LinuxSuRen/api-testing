# A wrapper to test.
#
# All make targets related to e2e,full test are defined in this file.

.PHONY: test.e2e
test.e2e:
	cd e2e && ./start.sh && ./start.sh compose-k8s.yaml && ./start.sh compose-external.yaml

.PHONY: test.fuzz
test.fuzz:
	cd pkg/util && go test -fuzz FuzzZeroThenDefault -fuzztime 6s

##@ Other Test

.PHONY: test-e2e
test-e2e: ## Run e2e tests
test-e2e: test.e2e

.PHONY: full-test-e2e
full-test-e2e: ## Build image then run e2e tests
full-test-e2e:
	TAG=master REGISTRY=ghcr.io make image test-e2e

.PHONY: test-fuzz
test-fuzz: ## Run fuzz tests
test-fuzz: test.fuzz
