#!/usr/bin/env bash
set -euo pipefail
echo "Creating Kafka topics..."
docker compose exec -T kafka kafka-topics.sh --create --if-not-exists \
  --bootstrap-server kafka:9092 --replication-factor 1 --partitions 3 \
  --topic trips_raw || true
echo "Topics created."
