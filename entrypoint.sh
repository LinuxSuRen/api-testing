#!/bin/sh

k3d cluster create
k3d cluster list

atest init -k "$2" --wait-namespace "$3" --wait-resource "$4"
atest run -p "$1"
