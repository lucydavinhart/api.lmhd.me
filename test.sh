#!/bin/bash

set -e

API_PATH="https://api.test.lmhd.me"

curl ${API_PATH}/v1/name
curl ${API_PATH}/v1/front

curl ${API_PATH}/v1/hello