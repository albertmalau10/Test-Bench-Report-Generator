#!/usr/bin/env bash
# For ctrlX CORE virtual only.
set -e
cd "$(dirname "$0")"
./build-frontend.sh
snapcraft clean
snapcraft pack --build-for=amd64 --verbosity=verbose
