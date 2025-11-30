#!/bin/sh
set -e

# Set defaults if not provided
export VITE_BACKEND_API="${VITE_BACKEND_API:-http://localhost:8081/api}"
export VITE_DEMO_MODE="${VITE_DEMO_MODE:-false}"

echo "Generating runtime-config.json with:"
echo "  VITE_BACKEND_API=${VITE_BACKEND_API}"
echo "  VITE_DEMO_MODE=${VITE_DEMO_MODE}"

# Generate runtime-config.json from template
envsubst < /usr/share/nginx/html/runtime-config.template.json > /usr/share/nginx/html/runtime-config.json

echo "Runtime configuration generated successfully."

# Start nginx
exec nginx -c /etc/nginx/nginx.conf -g "daemon off;"
