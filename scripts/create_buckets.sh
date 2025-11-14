#!/usr/bin/env bash
set -euo pipefail
echo "Creating MinIO buckets..."
docker compose exec -T minio bash -lc 'mc alias set local http://localhost:9000 minio minio12345 || true; mc mb -p local/mobility-raw || true; mc mb -p local/mobility-delta || true; mc mb -p local/mobility-marts || true'
echo "Done."
