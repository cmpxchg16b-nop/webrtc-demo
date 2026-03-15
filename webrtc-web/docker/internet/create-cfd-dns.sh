#!/bin/bash

# Create DNS record for cloudflared proxied domain, call this after ./create-tunnel.sh

set -e

script_path=$(realpath $0)
script_dir=$(dirname $script_path)
cd "${script_dir}"


source .env

if [ -z "${CF_API_TOKEN}" ]; then
    echo "No CF_API_TOKEN provided"
    exit 1
fi

if [ -z "${CF_ZONE_ID}" ]; then
    echo "No CF_ZONE_ID provided"
    exit 1
fi

if [ -z "${CLOUDFLARED_UUID}" ]; then
    echo "CLOUDFLARED_UUID is not set"
    exit 1
fi

if [ -z "${MAIN_DOMAIN}" ]; then
    echo "MAIN_DOMAIN is not set"
    exit 1
fi

curl "https://api.cloudflare.com/client/v4/zones/${CF_ZONE_ID}/dns_records" \
  --request POST \
  --header "Authorization: Bearer ${CF_API_TOKEN}" \
  --json "{
    \"type\": \"CNAME\",
    \"proxied\": true,
    \"name\": \"${MAIN_DOMAIN}\",
    \"content\": \"${CLOUDFLARED_UUID}.cfargotunnel.com\"
  }"
