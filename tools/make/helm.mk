# A wrapper to manage helm charts
#
# All make targets related to helmß are defined in this file.

include tools/make/env.mk


OCI_REGISTRY ?= oci://${REGISTRY}/${REGISTRY_NAMESPACE}
CHART_NAME ?= api-testing
CHART_VERSION ?= ${HELM_VERSION}

##@ Helm

.PHONY: helm-pkg
helm-pkg: ## Package API Testing helm chart.
helm-pkg: helm-dev-update
	@$(LOG_TARGET)
	# e.g. api-testing-v0.0.3-helm.tgz
	helm package helm/${CHART_NAME} --version ${CHART_VERSION}-helm --app-version ${CHART_VERSION} --destination ${OUTPUT_DIR}/charts/

.PHONY: helm-push
helm-push:
helm-push: ## Push API Testing helm chart to OCI registry.
	@$(LOG_TARGET)
	helm push ${OUTPUT_DIR}/charts/${CHART_NAME}-${CHART_VERSION}-helm.tgz ${OCI_REGISTRY}

.PHONY: helm-lint
helm-lint: ## Helm lint API Testing helm chart.
helm-lint: helm-dev-update
	helm lint helm/${CHART_NAME}

helm-dev-update:
helm-dev-update:
	helm dep update helm/${CHART_NAME}
