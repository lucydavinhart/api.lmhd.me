export VAULT_ADDR=https://vault.lmhd.me

# Issue client cert from Vault
vault write --format=json pki/inter/issue/api.test.lmhd.me-client common_name=issue-cert.sh > cert.json
cat cert.json | jq -r .data.certificate > cert.pem

# Inspect cert
cat cert.pem | openssl x509 -text -noout | head -n 11
