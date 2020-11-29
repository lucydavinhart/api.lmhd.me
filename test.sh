#!/bin/bash

set -e

if [[ -z ${API_PATH} ]]; then
	# Default to Test if unset
	API_PATH="https://api.test.lmhd.me"
fi

curl -s ${API_PATH}/hello | grep "HCL version"
curl -s ${API_PATH}/hello.hcl | grep "HCL version"
curl -s ${API_PATH}/hello.json | grep "JSON version"


# Check we can access it
curl ${API_PATH}/v1/name

# Check it shows name version
if [[ $(curl ${API_PATH}/v1/name.json | jq -r .version) == "1.2.0" ]]; then
	echo Name Contains Expected Version
else
	echo Unexpected name version $(curl ${API_PATH}/v1/name.json | jq .version)
	exit 1
fi

# Check we can access it
curl ${API_PATH}/v1/front
