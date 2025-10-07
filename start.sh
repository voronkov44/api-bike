#!/bin/bash
set -e

echo "Running migrations..."
/usr/local/bin/migrate

echo "Starting API..."
/usr/local/bin/app
