SHELL := /bin/bash
PROJECT := mobility-lakehouse

up:
\tdocker compose up -d --build

down:
\tdocker compose down -v

logs:
\tdocker compose logs -f --tail=200

init:
\tbash scripts/init_kafka_topics.sh
\tbash scripts/create_buckets.sh

airflow-sh:
\tdocker compose exec airflow-webserver bash

spark-submit:
\tdocker compose run --rm spark-submit \
\t  /opt/spark/bin/spark-submit \
\t  --class com.mobility.BronzeStream \
\t  --packages org.apache.spark:spark-sql-kafka-0-10_2.12:3.5.0,io.delta:delta-core_2.12:2.4.0,com.amazonaws:aws-java-sdk-bundle:1.12.743,org.apache.hadoop:hadoop-aws:3.3.4 \
\t  /jobs/mobility-stream-assembly.jar

dbt:
\tdocker compose run --rm dbt-run

ge:
\tdocker compose run --rm great-exp bash -lc "great_expectations checkpoint run bronze_checkpoint"
