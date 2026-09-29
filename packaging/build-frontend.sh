#!/usr/bin/env bash
# Builds the Vue frontend for the ctrlX CORE snap (base path
# "/valve-database-app/", matching main.go's proxyPrefix and
# package-manifest.json's proxyMapping url) and copies the result
# into packaging/www — read by snapcraft.yaml's "frontend" part.
# Run this BEFORE build-snap-amd64.sh / build-snap-arm64.sh.
set -e
cd "$(dirname "$0")"

rm -rf www
mkdir -p www

(
  cd ../frontend
  npm install
  VITE_BASE_PATH=/valve-database-app/ npm run build
)

cp -r ../frontend/dist/. www/
echo "frontend built into packaging/www/"
