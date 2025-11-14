# Project Playbook

## Run Order
1. `make up` — boot services
2. `make init` — buckets, topics
3. Confirm **Kafka** has topic `trips_raw`
4. Rust **trip-sim** is publishing (check `docker compose logs rust-trip-sim`)
5. `make spark-submit` — start Bronze streaming (or deploy to Airflow later)
6. Visit **MinIO**; observe objects in `mobility-delta/bronze/trips/...`
7. In **Grafana**, open "Lakehouse Freshness"
8. In **Trino**, run `sql/queries.sql`
9. In **Airflow**, enable `mobility_lakehouse` DAG (GTFS ingest → GE → dbt)
10. Query **gold** KPIs in Trino and visualize externally (e.g. Power BI).

## Resume Bullets (paste-ready)
- Delivered a **polyglot mobility lakehouse**: **Scala Spark** streaming from Kafka to **Delta Bronze**, **Python/Airflow** orchestration with **Great Expectations** data quality, **Ruby** GTFS ingestion to S3 (MinIO), and **Trino/dbt** curated marts; implemented **freshness SLAs** and **partition‑pruning** performance baselines.
- Built **Rust** load generator publishing ride events to Kafka; added a **Go** partition watcher exposing **Prometheus** metrics visualized in **Grafana**; automated bucket/topic init and one-command developer bring-up via Docker Compose and Make.

## Power BI Ideas
- Trips per hour by market (stacked column) + 7‑day trendline
- Avg trip duration vs distance scatter with market color
- Completion metrics: completed vs canceled counts
- Fare distribution histogram; top/bottom deciles
- Freshness indicator (last modified vs current time)
- (Optional) Map visual of pickup density if you later add geo tiles

## Notes
- All images and tools are OSS. You may swap versions as needed.
- Delta + Trino requires a functional Hive Metastore (provided).
- For Slack alerts from the Go watcher, set `SLACK_WEBHOOK_URL` env.
