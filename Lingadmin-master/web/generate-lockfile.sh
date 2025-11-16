#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(dirname "$0")
echo "Generating package-lock.json in $ROOT_DIR..."
if ! command -v npm >/dev/null 2>&1; then
  echo "npm is not installed. Please install Node.js/npm and retry." >&2
  exit 1
fi

cd "$ROOT_DIR"
# Generate lockfile without installing node_modules (supported by npm)
if npm --version >/dev/null 2>&1; then
  echo "Running: npm install --package-lock-only"
  npm install --package-lock-only
  echo "package-lock.json generated at $ROOT_DIR/package-lock.json"
else
  echo "npm not available to generate lockfile" >&2
  exit 1
fi
