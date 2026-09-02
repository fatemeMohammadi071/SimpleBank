#!/bin/sh
# ↑↑
# #! = shebang — "hey OS, use the following program to run me"
# /bin/sh = the path to the shell program

set -e

# DB migrations now run inside the Go binary at startup (see runDBMigration in main.go).

echo "starting server..."
exec "$@"