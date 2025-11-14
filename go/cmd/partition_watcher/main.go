package main

import (
\t"context"
\t"log"
\t"net/http"
\t"os"
\t"time"

\t"github.com/minio/minio-go/v7"
\t"github.com/minio/minio-go/v7/pkg/credentials"
\t"github.com/prometheus/client_golang/prometheus"
\t"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
\tobjCount = prometheus.NewGaugeVec(
\t\tprometheus.GaugeOpts{Name: "s3_partition_object_count", Help: "Objects per prefix"},
\t\t[]string{"prefix"},
\t)
\tlastMod = prometheus.NewGaugeVec(
\t\tprometheus.GaugeOpts{Name: "s3_partition_last_modified_epoch", Help: "Last modified time per prefix"},
\t\t[]string{"prefix"},
\t)
)

func main() {
\tprometheus.MustRegister(objCount, lastMod)

\tendpoint := getenv("MINIO_ENDPOINT", "http://minio:9000")
\tak := getenv("MINIO_ACCESS_KEY", "minio")
\tsk := getenv("MINIO_SECRET_KEY", "minio12345")
\tbucket := getenv("MINIO_BUCKET", "mobility-delta")
\tprefix := getenv("WATCH_PREFIX", "bronze/trips")

\tclient, err := minio.New(endpointHost(endpoint), &minio.Options{
\t\tCreds:  credentials.NewStaticV4(ak, sk, ""),
\t\tSecure: false,
\t})
\tif err != nil { log.Fatal(err) }

\tgo func() {
\t\tfor {
\t\t\tctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
\t\t\tcount := 0
\t\t\tvar latest time.Time
\t\t\tfor obj := range client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
\t\t\t\tif obj.Err != nil { continue }
\t\t\t\tcount++
\t\t\t\tif obj.LastModified.After(latest) { latest = obj.LastModified }
\t\t\t}
\t\t\tobjCount.WithLabelValues(prefix).Set(float64(count))
\t\t\tif !latest.IsZero() { lastMod.WithLabelValues(prefix).Set(float64(latest.Unix())) }
\t\t\tcancel()
\t\t\ttime.Sleep(15 * time.Second)
\t\t}
\t}()

\thttp.Handle("/metrics", promhttp.Handler())
\tlog.Println("Serving metrics on :9108/metrics")
\tlog.Fatal(http.ListenAndServe(":9108", nil))
}

func getenv(k, d string) string {
\tif v := os.Getenv(k); v != "" { return v }
\treturn d
}
func endpointHost(e string) string {
\t// strip scheme if present
\tif len(e) > 7 && e[:7] == "http://" { return e[7:] }
\tif len(e) > 8 && e[:8] == "https://" { return e[8:] }
\treturn e
}
