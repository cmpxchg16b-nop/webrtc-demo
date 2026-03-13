#!/bin/bash

set -e

if [ -z $MANAGEMENT_API ]; then
    echo "MANAGEMENT_API is not set"
    exit 1
fi
MANAGEMENT_API=$(realpath $MANAGEMENT_API)
echo MANAGEMENT_API: $MANAGEMENT_API

if [ -z $AVATAR_B64_PATH ]; then
    echo "AVATAR_B64_PATH is not set"
    exit 1
fi
AVATAR_B64_PATH=$(realpath $AVATAR_B64_PATH)
echo AVATAR_B64_PATH: $AVATAR_B64_PATH

current_path=$(realpath $0)
cd $(dirname $current_path)/..

echo Using $(pwd) as current working directory

# invocation:
# scripts/launch-test-agents-server.sh
# and .env should be in $PWD

MANAGEMENT_API=${MANAGEMENT_API} \
AVATAR_B64_PATH=${AVATAR_B64_PATH} \
scripts/prepare-bot-tokens.sh go run ./cmd/agents/main.go \
  --chatbot-model=deepseek/deepseek-v3.2
