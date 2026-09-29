#!/usr/bin/env bash
# For physical ctrlX CORE hardware (M3/M4 and similar).
set -e
cd "$(dirname "$0")"
./build-frontend.sh
snapcraft clean
snapcraft pack --build-for=arm64 --verbosity=verbose
