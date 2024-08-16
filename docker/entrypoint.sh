#!/bin/sh

set -e

echo "Running option $1"

# Check the first argument passed to the script
if [ "$1" = "-users-api" ]; then
  echo "Starting users-service" && ./users-service
  exit 0
fi


if [ "$1" = "-newsfeed-api" ]; then
  echo "Starting newsfeed-service" && ./newsfeed-service
  exit 0
fi