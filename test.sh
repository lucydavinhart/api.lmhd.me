#!/bin/bash

set -e
set -o pipefail

# https://superuser.com/a/1641410
# fail and output body
curlf() {
  OUTPUT_FILE=$(mktemp)
  HTTP_CODE=$(curl --silent --output $OUTPUT_FILE --write-out "%{http_code}" "$@")
  if [[ ${HTTP_CODE} -lt 200 || ${HTTP_CODE} -gt 299 ]] ; then
    >&2 cat $OUTPUT_FILE
    return 22
  fi
  cat $OUTPUT_FILE
  rm $OUTPUT_FILE
}

if [[ -z ${API_PATH} ]]; then
	# Default to Test if unset
	API_PATH="https://api.test.lmhd.me"
fi

echo
echo ========================================
echo Testing Hello World
echo ========================================
echo

curlf -s ${API_PATH}/hello | grep "HCL version"
curlf -s ${API_PATH}/hello.hcl | grep "HCL version"
curlf -s ${API_PATH}/hello.json | grep "JSON version"


echo
echo ========================================
echo Testing Name
echo ========================================
echo

# Check we can access it
curlf ${API_PATH}/v1/name

# Check it shows name version
if [[ $(curlf ${API_PATH}/v1/name.json | jq -r .version) == "1.2.0" ]]; then
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
curlf ${API_PATH}/v1/front



# TODO: Disable in Circle?
# Or give Circle the ability to generate Test certs?
echo
echo ========================================
echo Testing Auth
echo ========================================
echo

# Check we can access it
if [[ "${API_PATH}" == "https://api.lmhd.me" ]]; then
	./issue-cert.sh prod
else
	./issue-cert.sh
fi
curlf -H "X-LMHD-Client-Cert: $(base64 cert.pem)" ${API_PATH}/v1/auth


echo
echo ========================================
echo Testing Federate
echo ========================================
echo

# Check we can access it
if [[ "${API_PATH}" == "https://api.lmhd.me" ]]; then
	./issue-cert.sh prod
else
	./issue-cert.sh
fi
curlf -H "X-LMHD-Client-Cert: $(base64 cert.pem)" -X POST ${API_PATH}/v1/front/federate

