#!/bin/sh
# ↑↑
# #! = shebang — "hey OS, use the following program to run me"
# /bin/sh = the path to the shell program

set -e

echo "Run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "starting server..."
exec "$@"