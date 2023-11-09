#!/bin/bash

docker-compose version
docker-compose up --build -d

while true
do
    docker-compose ps | grep testing
    if [ $? -eq 1 ]
    then
        code=-1
        docker-compose logs | grep e2e-testing
        docker-compose logs | grep e2e-testing | grep Usage
        if [ $? -eq 1 ]
        then
            code=0
            echo "successed"
        fi

        docker-compose down
        set -e
        exit $code
    fi
    sleep 1
done
