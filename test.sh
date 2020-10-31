#!/bin/bash

set -ex

API_PATH="https://api.test.lmhd.me"

curl -s ${API_PATH}/hello | grep "HCL version"
curl -s ${API_PATH}/hello.hcl | grep "HCL version"
curl -s ${API_PATH}/hello.json | grep "JSON version"


curl ${API_PATH}/v1/name
curl ${API_PATH}/v1/front