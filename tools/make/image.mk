# A wrapper to build and push docker image
#
# All make targets related to docker image are defined in this file.

include tools/make/env.mk


# Determine image files by looking into ./Dockerfile
IMAGES_DIR ?= $(wildcard ${ROOT_DIR}tools/docker/*)
# Determine images names by stripping out the dir names
IMAGES ?= atest

ifeq (${IMAGES},)
  $(error Could not determine IMAGES, set ROOT_DIR or run in source dir)
endif

.PHONY: image.build
image.build: $(addprefix image.build., $(IMAGES))

.PHONY: image.build.%
# Maybe can use: image.build.%: go.build.$(GOOS)_$(GOARCH).%
image.build.%:
	@$(LOG_TARGET)
	@$(call log, "Building image $(GOOS)-$(GOARCH) $(IMAGES):$(TAG)")
	${IMAGE_TOOL} build -f $(ROOT_DIR)/Dockerfile \
			-t ${REGISTRY}/${IMAGES}:${TAG} . \
    		--build-arg GOPROXY=${GOPROXY} \
    		--build-arg VERSION=$(TAG)

.PHONY: run.image
run.image:
	@$(LOG_TARGET)
	${IMAGE_TOOL} run -p 7070:7070 -p 8080:8080 ${REGISTRY}/${IMAGES}:${TAG}

##@ Image

.PHONY: image
image: ## Build docker images for host platform. See Option PLATFORM and BINARY.
image: image.build

.PHONY: run-container
run-container: ## Run the docker container for the image. See Option IMAGES and TAG.
run-container: run.image
