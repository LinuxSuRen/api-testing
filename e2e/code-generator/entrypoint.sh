#!/bin/bash
set -e

export lang=$1
# exit if no language is specified
if [ -z "$lang" ]
then
  echo "no language is specified"
  exit 1
fi

mkdir -p /root/.config/atest
mkdir -p /var/data

nohup atest server --local-storage '/workspace/test-suites/*.yaml'&
sleep 1

test_cases=("requestWithHeader" "requestWithoutHeader")

for test_case in "${test_cases[@]}"
do
    curl http://localhost:8080/server.Runner/GenerateCode -X POST \
        -d '{"TestSuite": "test", "TestCase": "'"${test_case}"'", "Generator": "'"$lang"'"}' > code.json

    cat code.json | jq .message -r | sed 's/\\n/\n/g' | sed 's/\\t/\t/g' | sed 's/\\\"/"/g' > code.txt
    cat code.txt

    sh /workspace/${lang}.sh code.txt
done

exit 0
