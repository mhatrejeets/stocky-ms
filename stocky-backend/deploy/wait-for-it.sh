#!/bin/sh
# wait-for-it.sh: Wait until a host:port is available
# Usage: wait-for-it.sh host:port -- command args

set -e

hostport="$1"
shift

host="${hostport%%:*}"
port="${hostport##*:}"

while ! nc -z "$host" "$port"; do
  echo "Waiting for $host:$port..."
  sleep 2
done

echo "$host:$port is available. Starting app."
exec "$@"