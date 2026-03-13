#!/bin/bash

set -e

if [ -z $MANAGEMENT_API ]; then
    echo "MANAGEMENT_API is not set"
    exit 1
fi

if [ -z $AVATAR_B64_PATH ]; then
    echo "AVATAR_B64_PATH is not set"
    exit 1
fi

echo "Getting Bot JWT token for EchoBot ..."
echobot_token=$(MANAGEMENT_API=$MANAGEMENT_API scripts/create-bot-account.sh echo_bot EchoBot $AVATAR_B64_PATH)
if [ -z "${echobot_token}" ] || [ "${echobot_token}" == "null" ]; then
    echo "Failed to get echobot_token"
    exit 1
fi

echo "Getting Bot JWT token for MusicBot ..."
musicbot_token=$(MANAGEMENT_API=$MANAGEMENT_API scripts/create-bot-account.sh music_bot MusicBot $AVATAR_B64_PATH)
if [ -z "${musicbot_token}" ] || [ "${musicbot_token}" == "null" ]; then
    echo "Failed to get musicbot_token"
    exit 1
fi

echo "Getting Bot JWT token for CounterBot ..."
counterbot_token=$(MANAGEMENT_API=$MANAGEMENT_API scripts/create-bot-account.sh counter_bot CounterBot $AVATAR_B64_PATH)
if [ -z "${counterbot_token}" ] || [ "${counterbot_token}" == "null" ]; then
    echo "Failed to get counterbot_token"
    exit 1
fi

echo "Getting Bot JWT token for ClockBot ..."
clockbot_token=$(MANAGEMENT_API=$MANAGEMENT_API scripts/create-bot-account.sh clock_bot ClockBot $AVATAR_B64_PATH)
if [ -z "${clockbot_token}" ] || [ "${clockbot_token}" == "null" ]; then
    echo "Failed to get clockbot_token"
    exit 1
fi

echo "Getting Bot JWT token for ChatBot ..."
chatbot_token=$(MANAGEMENT_API=$MANAGEMENT_API scripts/create-bot-account.sh chat_bot ChatBot $AVATAR_B64_PATH)
if [ -z "${chatbot_token}" ] || [ "${chatbot_token}" == "null" ]; then
    echo "Failed to get chatbot_token"
    exit 1
fi

echo "Generating .env.bot file in current directory $PWD"
echo "ECHOBOT_JWT_TOKEN=$echobot_token" > .env.bots
echo "MUSICBOT_JWT_TOKEN=$musicbot_token" >> .env.bots
echo "COUNTERBOT_JWT_TOKEN=$counterbot_token" >> .env.bots
echo "CLOCKBOT_JWT_TOKEN=$clockbot_token" >> .env.bots
echo "CHATBOT_JWT_TOKEN=$chatbot_token" >> .env.bots

echo "Done, launching bot agents process ..."
exec "$@"
