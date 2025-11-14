# Mobility Lakehouse — Polyglot (Spark/Scala, Python/Airflow, Ruby, Go, Rust, SQL/Trino, TypeScript UI)

A **free, local, dockerized** end‑to‑end project demonstrating data engineering skills across the stack:

- **Scala (Spark Structured Streaming)**: Kafka → **Delta Bronze**, aggregates → **Silver/Gold**
- **Python (Airflow + Great Expectations)**: Orchestration, data quality checks, SLAs
- **Ruby (CLI)**: GTFS ingest → MinIO (S3), partitioning + validation
- **Go (service)**: Partition watcher → **Prometheus** metrics (+ optional Slack webhook)
- **Rust (load generator)**: Trip simulator publishing JSON events to Kafka
- **SQL (Trino)**: Marts, query tuning (partition pruning, joins), dbt models + docs
- **TypeScript UI (React)**: Micro data catalog & sample queries
- **Grafana + Prometheus**: Charts of pipeline freshness and volumes

All services run **locally** via Docker Compose. No paid services required.

---

## Quick Start

**Prereqs:** Docker Desktop or Docker Engine + Compose v2; Git; Node 18+ (optional for UI), Java (optional for local Spark build).

```bash
git clone <your-fork-url> mobility-lakehouse
cd mobility-lakehouse
make up     # boots the full stack
```

Then open:
- MinIO: http://localhost:9001  (user: `minio`, pass: `minio12345`)
- Airflow: http://localhost:8080 (user: `airflow`, pass: `airflow`)
- Trino Web UI: http://localhost:8082
- Grafana: http://localhost:3000 (user: `admin`, pass: `admin`)
- Prometheus: http://localhost:9090
- Data UI (React): http://localhost:5173

### One-time setup
```bash
make init   # creates S3 buckets, Kafka topics, and Hive schemas
```

### Run the pipeline
- The **Rust trip simulator** auto-publishes JSON to Kafka (`trips_raw`).
- The **Spark streaming job** (Scala) consumes Kafka and writes **Delta Bronze** to MinIO.
- The **Airflow DAG** runs **Ruby GTFS ingest**, **Great Expectations** checks, and **dbt** to produce **Silver/Gold** marts.
- **Go watcher** exposes freshness metrics (Prometheus), **Grafana** shows dashboards.
- Query with **Trino** (Delta + Hive/Parquet): see `sql/` and the **UI** for examples.

---

## Services & Ports

| Service         | Port  | Notes |
|-----------------|-------|------|
| Zookeeper       | 2181  | Kafka dependency |
| Kafka           | 9092  | PLAINTEXT broker |
| MinIO           | 9000  | S3 API |
| MinIO Console   | 9001  | Web console |
| Postgres        | 5432  | Metastore DB + Airflow DB |
| Hive Metastore  | 9083  | Metastore for Trino/Spark |
| Trino           | 8082  | Web UI + JDBC |
| Airflow Web     | 8080  | Orchestration |
| Prometheus      | 9090  | Metrics |
| Grafana         | 3000  | Dashboards |
| Go Watcher      | 9108  | /metrics |
| React UI        | 5173  | Data catalog |

---

## Data Sources

- **Trips (simulated)**: Rust generator publishes JSON to Kafka topic `trips_raw` at ~X msgs/sec.
- **GTFS (real public data)**: Ruby CLI downloads a GTFS ZIP (configurable URL) and lands partitioned CSV to `s3://mobility-raw/gtfs/YYYY/MM/DD/`.

> Default GTFS URL is set to Calgary Public Transit (can be changed via `GTFS_URL` env).

---

## Useful Commands

```bash
make up            # start everything
make down          # stop & remove containers
make init          # init buckets/topics/schemas
make logs          # tail compose logs
make airflow-sh    # shell into airflow worker
make spark-submit  # run the Scala Spark job (bronze streaming)
make dbt           # run dbt models (silver/gold)
make ge            # run Great Expectations checkpoint
```

## Next Steps

- Open the **UI** at http://localhost:5173 to browse catalogs/schemas/tables.
- In **Trino**, run the queries in `sql/queries.sql` (partition pruning demos).
- Explore **Grafana** dashboard (freshness, object counts).

See `docs/PLAYBOOK.md` for a full walkthrough and resume bullets.
