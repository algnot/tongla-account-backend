#!/bin/sh

echo "🚀 Running migration..."
go run ./cmd --with-migrate

echo "✅ Starting app..."
go run ./cmd
