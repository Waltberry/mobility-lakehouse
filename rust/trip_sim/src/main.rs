use chrono::{Utc};
use rand::Rng;
use rdkafka::config::ClientConfig;
use rdkafka::producer::{FutureProducer, FutureRecord};
use serde::Serialize;
use std::time::Duration;
use std::{env, thread};

#[derive(Serialize)]
struct Trip {
    trip_id: String,
    rider_id: String,
    driver_id: String,
    market: String,
    status: String,
    ts: String,
    pickup_lat: f64,
    pickup_lon: f64,
    dropoff_lat: f64,
    dropoff_lon: f64,
    duration_sec: i32,
    distance_km: f64,
    fare: f64,
}

fn main() {
    let broker = env::var("KAFKA_BROKER").unwrap_or("kafka:9092".into());
    let topic = env::var("TOPIC").unwrap_or("trips_raw".into());
    let mps: u64 = env::var("MSG_PER_SEC").unwrap_or("10".into()).parse().unwrap_or(10);

    let producer: FutureProducer = ClientConfig::new()
        .set("bootstrap.servers", &broker)
        .create()
        .expect("producer");

    let markets = vec!["CA-AB", "CA-ON", "CA-BC"];
    let mut rng = rand::thread_rng();

    loop {
        for _ in 0..mps {
            let trip = Trip {
                trip_id: format!("t{}", rng.gen::<u64>()),
                rider_id: format!("r{}", rng.gen::<u32>()),
                driver_id: format!("d{}", rng.gen::<u32>()),
                market: markets[rng.gen_range(0..markets.len())].to_string(),
                status: "completed".into(),
                ts: Utc::now().to_rfc3339(),
                pickup_lat: 51.0 + rng.gen::<f64>() * 0.1,
                pickup_lon: -114.0 + rng.gen::<f64>() * 0.1,
                dropoff_lat: 51.0 + rng.gen::<f64>() * 0.1,
                dropoff_lon: -114.0 + rng.gen::<f64>() * 0.1,
                duration_sec: rng.gen_range(300..3600),
                distance_km: rng.gen_range(1.0..35.0),
                fare: rng.gen_range(5.0..80.0),
            };
            let payload = serde_json::to_string(&trip).unwrap();
            let record = FutureRecord::to(&topic).payload(&payload).key(&trip.trip_id);
            let _ = producer.send(record, Duration::from_secs(0));
        }
        thread::sleep(Duration::from_secs(1));
    }
}
