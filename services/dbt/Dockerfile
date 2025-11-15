FROM python:3.11-slim

# system deps
RUN apt-get update && apt-get install -y --no-install-recommends \
    git build-essential libpq-dev curl ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# dbt + trino adapter
RUN pip install --no-cache-dir dbt-core==1.8.6 dbt-trino==1.8.4 pandas pyarrow

WORKDIR /usr/app
ENTRYPOINT ["dbt"]
