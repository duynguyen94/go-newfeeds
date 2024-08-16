#!/bin/sh

set -e

until nc -z -v -w30 "$MYSQL_HOST" "$MYSQL_PORT"
do
  echo "Waiting for database connection..."
  sleep 1
done

echo "MySQL is up and running"

until nc -z -v -w30 "$REDIS_HOST" "$REDIS_PORT"
do
  echo "Waiting for redis connection..."
  sleep 1
done

echo "Redis is up and running"

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