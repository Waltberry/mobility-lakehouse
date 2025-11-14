package com.mobility

import org.apache.spark.sql.SparkSession
import org.apache.spark.sql.functions._
import org.apache.spark.sql.types._

object BronzeStream {
  def main(args: Array[String]): Unit = {
    val spark = SparkSession.builder()
      .appName("Mobility Bronze Stream")
      .getOrCreate()

    import spark.implicits._

    val schema = new StructType()
      .add("trip_id", StringType)
      .add("rider_id", StringType)
      .add("driver_id", StringType)
      .add("market", StringType)
      .add("status", StringType)
      .add("ts", TimestampType)
      .add("pickup_lat", DoubleType)
      .add("pickup_lon", DoubleType)
      .add("dropoff_lat", DoubleType)
      .add("dropoff_lon", DoubleType)
      .add("duration_sec", IntegerType)
      .add("distance_km", DoubleType)
      .add("fare", DoubleType)

    val kafkaBootstrap = sys.env.getOrElse("KAFKA_BOOTSTRAP", "kafka:9092")
    val topic = sys.env.getOrElse("KAFKA_TOPIC", "trips_raw")
    val s3Path = sys.env.getOrElse("DELTA_S3_PATH", "s3a://mobility-delta/bronze/trips")
    val checkpointLoc = sys.env.getOrElse("CHECKPOINT_PATH", "s3a://mobility-delta/_checkpoints/bronze_trips")

    spark.sparkContext.hadoopConfiguration.set("fs.s3a.endpoint", sys.env.getOrElse("S3_ENDPOINT","http://minio:9000"))
    spark.sparkContext.hadoopConfiguration.set("fs.s3a.path.style.access", "true")
    spark.sparkContext.hadoopConfiguration.set("fs.s3a.access.key", sys.env.getOrElse("AWS_ACCESS_KEY_ID","minio"))
    spark.sparkContext.hadoopConfiguration.set("fs.s3a.secret.key", sys.env.getOrElse("AWS_SECRET_ACCESS_KEY","minio12345"))
    spark.sparkContext.hadoopConfiguration.set("fs.s3a.connection.ssl.enabled", "false")

    val df = spark.readStream
      .format("kafka")
      .option("kafka.bootstrap.servers", kafkaBootstrap)
      .option("subscribe", topic)
      .option("startingOffsets", "latest")
      .load()
      .selectExpr("CAST(value AS STRING) AS json_str")
      .select(from_json($"json_str", schema).as("r"))
      .select("r.*")
      .withColumn("event_date", to_date($"ts"))

    val query = df.writeStream
      .format("delta")
      .outputMode("append")
      .partitionBy("event_date","market")
      .option("checkpointLocation", checkpointLoc)
      .start(s3Path)

    query.awaitTermination()
  }
}
