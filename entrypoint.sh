#!/bin/sh
set -e

echo "🚀 Running migration..."
./main --with-migrate

echo "✅ Starting app..."
exec ./main
