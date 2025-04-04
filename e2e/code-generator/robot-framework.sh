#!/bin/bash
set -e

export sourcefile=$1
# exit if no source file is specified
if [ -z "$sourcefile" ]
then
  echo "no source file is specified"
  exit 1
fi

mv ${sourcefile} test.robot
pip install robotframework robotframework-requests
robot test.robot
