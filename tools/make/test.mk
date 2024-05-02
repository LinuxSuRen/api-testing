# A wrapper to test.
#
# All make targets related to e2e,operator,full test are defined in this file.

.PHONY: test.operator
test.operator:
	cd operator && \
	go test ./... -cover -v -coverprofile=coverage.out && \
	go tool cover -func=coverage.out

.PHONY: test.e2e
test.e2e:
	cd e2e && ./start.sh && ./start.sh compose-k8s.yaml && ./start.sh compose-external.yaml

.PHONY: test.fuzz
test.fuzz:
	cd pkg/util && go test -fuzz FuzzZeroThenDefault -fuzztime 6s

##@ Other Test

.PHONY: test-operator
test-operator: ## Run operator tests
test-operator: test.operator

.PHONY: test-e2e
test-e2e: ## Run e2e tests
test-e2e: test.e2e

.PHONY: test-fuzz
test-fuzz: ## Run fuzz tests
test-fuzz: test.fuzz
