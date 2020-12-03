#!/bin/bash

set -e
set -o pipefail


if [[ -z ${API_PATH} ]]; then
	# Default to Test if unset
	API_PATH="https://api.test.lmhd.me"
fi

echo
echo ========================================
echo Testing Hello World
echo ========================================
echo

curl --fail -s ${API_PATH}/hello | grep "HCL version"
curl --fail -s ${API_PATH}/hello.hcl | grep "HCL version"
curl --fail -s ${API_PATH}/hello.json | grep "JSON version"


echo
echo ========================================
echo Testing Name
echo ========================================
echo

# Check we can access it
curl --fail ${API_PATH}/v1/name

# Check it shows name version
if [[ $(curl --fail ${API_PATH}/v1/name.json | jq -r .version) == "1.2.0" ]]; then
	echo Name Contains Expected Version
else
	echo Unexpected name version $(curl ${API_PATH}/v1/name.json | jq .version)
	exit 1
fi

echo
echo ========================================
echo Testing Fronter
echo ========================================
echo

# Check we can access it
curl --fail ${API_PATH}/v1/front
