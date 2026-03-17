#!/bin/bash

set -e

script_path=$(realpath $0)
script_dir=$(dirname $script_path)

cd $script_dir

source .env

if [ -z "${CADDY_ADMIN_ADDR}" ]; then
    echo "Can't find caddy admin addr"
    exit 1
fi

curl \
 -H 'Content-Type: application/json' \
 -H 'Cache-Control: must-revalidate' \
 -X POST \
 -d @caddy.apps.http.json \
 "http://${CADDY_ADMIN_ADDR}/config/apps/http"
