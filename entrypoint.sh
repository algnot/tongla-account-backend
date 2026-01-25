#!/bin/sh
set -e

if [ "$RUN_MIGRATE" = "true" ]; then
  echo "🚀 Running migration..."
  ./main --with-migrate
  echo "✅ Migration done"
  exit 0
fi

echo "✅ Starting app..."
exec ./main
