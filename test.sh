#!/bin/bash

set -ex

if [[ -z ${API_PATH} ]]; then
	# Default to Test if unset
	API_PATH="https://api.test.lmhd.me"
fi

curl -s ${API_PATH}/hello | grep "HCL version"
curl -s ${API_PATH}/hello.hcl | grep "HCL version"
curl -s ${API_PATH}/hello.json | grep "JSON version"


curl ${API_PATH}/v1/name
curl ${API_PATH}/v1/front
