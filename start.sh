#!/bin/sh

set -e
# FILE=/app/.env
# if [ ! -f "$FILE" ]; then
#     touch .env
# fi



echo "start up"
exec "$@"