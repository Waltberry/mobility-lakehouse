from datetime import datetime, timedelta
from airflow import DAG
from airflow.operators.bash import BashOperator

default_args = {
    "owner": "airflow",
    "retries": 1,
    "retry_delay": timedelta(minutes=2),
}

with DAG(
    dag_id="mobility_lakehouse",
    start_date=datetime(2025, 1, 1),
    schedule_interval="*/15 * * * *",
    catchup=False,
    default_args=default_args,
    tags=["lakehouse","mobility"],
) as dag:

    ingest_gtfs = BashOperator(
        task_id="ingest_gtfs",
        bash_command="ruby /workspace/ruby/bin/ingest_gtfs --region CA-AB",
        env={"GTFS_URL": "https://data.calgary.ca/download/npk7-z3bj/application/zip", "S3_BUCKET":"mobility-raw"}
    )

    spark_bronze = BashOperator(
        task_id="spark_bronze",
        bash_command="python -c 'print("Spark job triggered via Makefile; run `make spark-submit` manually in dev")'"
    )

    ge_bronze = BashOperator(
        task_id="great_expectations_bronze",
        bash_command="great_expectations checkpoint run bronze_checkpoint",
        cwd="/opt/airflow/great_expectations"
    )

    dbt_run = BashOperator(
        task_id="dbt_run",
        bash_command="dbt run",
        cwd="/usr/app/mobility"
    )

    ingest_gtfs >> spark_bronze >> ge_bronze >> dbt_run
