FROM rust:1.77 as build
RUN apt-get update && apt-get install -y cmake pkg-config libssl-dev libsasl2-dev
WORKDIR /app
COPY . .
RUN cargo build --release

FROM debian:12-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /
COPY --from=build /app/target/release/trip_sim /trip_sim
ENV KAFKA_BROKER=kafka:9092 TOPIC=trips_raw MSG_PER_SEC=10
ENTRYPOINT ["/trip_sim"]
