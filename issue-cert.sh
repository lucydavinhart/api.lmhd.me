export VAULT_ADDR=https://vault.lmhd.me

pki_role=api.test.lmhd.me-client
if [[ "$1" == "prod" ]]; then
	pki_role=api.lmhd.me-client
fi

# Issue client cert from Vault
vault write --format=json pki/api/issue/${pki_role} common_name=issue-cert.sh > cert.json
cat cert.json | jq -r .data.certificate > cert.pem

# Inspect cert
cat cert.pem | openssl x509 -text -noout | head -n 11
