# rust/trip_sim/Dockerfile
FROM rust:1.84 AS build

WORKDIR /app

# Install native deps needed by rdkafka (librdkafka)
RUN apt-get update && apt-get install -y --no-install-recommends \
    cmake pkg-config libssl-dev libsasl2-dev ca-certificates \
 && rm -rf /var/lib/apt/lists/*

# Copy only the manifest first and create a lockfile (works even if none exists)
COPY Cargo.toml ./
RUN cargo generate-lockfile

# Pre-cache dependencies (dummy main)
RUN mkdir -p src && echo "fn main(){}" > src/main.rs \
 && cargo fetch

# Now bring in your real code and build
COPY src ./src
RUN cargo build --release

# --------- runtime stage ----------
FROM debian:12-slim
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/target/release/trip_sim /app/trip_sim
ENV RUST_LOG=info
ENTRYPOINT ["/app/trip_sim"]
