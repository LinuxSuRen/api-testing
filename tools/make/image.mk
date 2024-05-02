# A wrapper to build and push docker image
#
# All make targets related to docker image are defined in this file.

include tools/make/env.mk


# Determine image files by looking into ./Dockerfile
IMAGES_DIR ?= ${ROOT_PACKAGE}/Dockerfile
# Determine images names by stripping out the dir names
IMAGES ?= atest

ifeq (${IMAGES},)
  $(error Could not determine IMAGES, set ROOT_DIR or run in source dir)
endif

.PHONY: image.build
image.build: $(addprefix image.build., $(IMAGES))

.PHONY: image.build.%
image.build.%: go.build.linux_$(GOARCH).%
	@$(LOG_TARGET)
	$(eval COMMAND := $(word 1,$(subst ., ,$*)))
	$(eval IMAGES := $(COMMAND))
	@$(call log, "Building image $(IMAGES):$(TAG)")
	${IMG_TOOL} build -t ${REGISTRY}/${IMAGES}:${TAG} . \
    		--build-arg GOPROXY=${GOPROXY} \
    		--build-arg VERSION=$(shell git describe --abbrev=0 --tags)-$(TAG)
	@$(call log, "Creating image tag $(REGISTRY)/$(IMAGES):$(TAG)")

.PHONY: run.image
run.image:
	@$(LOG_TARGET)
	${IMG_TOOL} run -p 7070:7070 -p 8080:8080 ${REGISTRY}/${IMAGES}:${TAG}

##@ Image

.PHONY: image
image: ## Build docker images for host platform. See Option PLATFORM and BINS.
image: image.build

.PHONY: run-container
run-container: ## Run the docker container for the image. See Option IMAGES and TAG.
run-container: run.image
