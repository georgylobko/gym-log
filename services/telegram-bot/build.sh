#!/usr/bin/env bash
# Builds the Lambda deployment directory: installs runtime deps + copies source
# into dist/. CDK ships dist/ as the function code (no Docker required).
#
# Builds into a temp dir and only swaps it into place on success, so a failed or
# interrupted build can never leave a half-populated dist/.
#
# Run from anywhere; paths are resolved relative to this script.
#   ./build.sh
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST="$HERE/dist"
BUILD="$(mktemp -d "${TMPDIR:-/tmp}/gymbot-build.XXXXXX")"

# Clean up the temp dir on any exit (success or failure).
trap 'rm -rf "$BUILD"' EXIT

echo "Installing dependencies into staging dir"
# --platform/--only-binary keeps wheels compatible with the Lambda runtime even
# when building on macOS. Pure-Python deps are unaffected.
python3 -m pip install \
  --target "$BUILD" \
  --requirement "$HERE/requirements.txt" \
  --platform manylinux2014_x86_64 \
  --python-version 3.12 \
  --implementation cp \
  --only-binary=:all: \
  --upgrade

echo "Copying source"
cp -r "$HERE/src/." "$BUILD/"

# Verify the build produced a usable package BEFORE touching dist/.
if [[ ! -f "$BUILD/lambda_function.py" ]]; then
  echo "ERROR: build is missing lambda_function.py — not publishing." >&2
  exit 1
fi

echo "Publishing to $DIST"
rm -rf "$DIST"
mv "$BUILD" "$DIST"

echo "Build complete: $DIST"
