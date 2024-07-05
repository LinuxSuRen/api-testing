#!/bin/bash
set -e

docker compose version

targets=(golang java python javascript curl)
for target in "${targets[@]}"
do
    docker compose down
    docker compose up --build $target --exit-code-from $target
done
