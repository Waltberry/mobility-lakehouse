name := "mobility-stream"
version := "0.1.0"
scalaVersion := "2.12.18"

libraryDependencies ++= Seq(
  "org.apache.spark" %% "spark-sql" % "3.5.0" % "provided",
  "org.apache.spark" %% "spark-sql-kafka-0-10" % "3.5.0",
  "io.delta" %% "delta-core" % "2.4.0",
  "com.amazonaws" % "aws-java-sdk-bundle" % "1.12.743",
  "org.apache.hadoop" % "hadoop-aws" % "3.3.4"
)
