#!/bin/bash

# Exit if any command fails
set -e

echo "Starting EMQX TLS certificate generation..."

# Generate Root CA key and certificate
openssl genpkey -algorithm RSA -out rootCA.key
openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 3650 -out rootCA.crt -subj "/CN=My-MQTT-CA"

# Generate Server key and Certificate Signing Request (CSR)
openssl genpkey -algorithm RSA -out server.key
openssl req -new -key server.key -out server.csr -subj "/CN=emqx"

# Create a temporary OpenSSL configuration file for the Subject Alternative Name (SAN)
# The `alt_names` section must be correctly formatted to be read by openssl.
cat > openssl.cnf <<-EOF
[req]
distinguished_name = req_distinguished_name
[req_distinguished_name]
[v3_ext]
subjectAltName = DNS:emqx
EOF

# Sign the server certificate with the Root CA using the new config file
openssl x509 -req -in server.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out server.crt -days 365 -sha256 -extfile openssl.cnf -extensions v3_ext

echo "Certificate generation complete. Files created:"
ls