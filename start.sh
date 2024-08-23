#!/bin/sh

echo "run db migration"

if [ ! -f /app/app.env ]; then
  echo "app.env not found!"
  exit 1
fi

set -a
. /app/app.env
set +a

grep -v '^#' app.env | grep -v '^$'

# Debugging: Print the DB_SOURCE to ensure it is loaded
echo "DB_SOURCE is $DB_SOURCE"


/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"

exec "$@" 