# A wrapper to hold common environment variables used in other make wrappers
#
# This file does not contain any specific make targets.

# Docker variables

HELM_VERSION ?= v0.0.3

# Docker repo
REGISTRY ?= docker.io/linuxsuren

# Set image tools
IMAGE_TOOL ?= docker

TOOLEXEC ?= #-toolexec="skywalking-go-agent"

GOPROXY ?= direct

TAG ?= ${REV}

FRONT_RUNTIMES ?= npm
