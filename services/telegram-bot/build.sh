#!/usr/bin/env bash
# Builds the Lambda deployment directory: installs runtime deps + copies source
# into dist/. CDK ships dist/ as the function code (no Docker required).
#
# Run from anywhere; paths are resolved relative to this script.
#   ./build.sh
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST="$HERE/dist"

echo "Cleaning $DIST"
rm -rf "$DIST"
mkdir -p "$DIST"

echo "Installing dependencies into $DIST"
# --platform/--only-binary keeps wheels compatible with the Lambda runtime even
# when building on macOS. Pure-Python deps are unaffected.
python3 -m pip install \
  --target "$DIST" \
  --requirement "$HERE/requirements.txt" \
  --platform manylinux2014_x86_64 \
  --python-version 3.12 \
  --implementation cp \
  --only-binary=:all: \
  --upgrade

echo "Copying source"
cp -r "$HERE/src/." "$DIST/"

echo "Build complete: $DIST"
