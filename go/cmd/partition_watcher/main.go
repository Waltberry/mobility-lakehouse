package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricObjCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "lakehouse_objects_total",
		Help: "Count of objects under the watched prefix",
	})
	metricLastMod = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "lakehouse_last_modified_seconds",
		Help: "Unix timestamp of the most-recent object",
	})
	metricFreshness = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "lakehouse_freshness_seconds",
		Help: "Seconds since last object modification",
	})
)

func mustEnv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func parseSecure(endpoint string) (host string, secure bool) {
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme == "" {
		// assume http if scheme missing
		return endpoint, false
	}
	host = u.Host
	secure = (u.Scheme == "https")
	return
}

func postSlack(webhook, text string) {
	body := map[string]string{"text": text}
	b, _ := json.Marshal(body)
	_, err := http.Post(webhook, "application/json", strings.NewReader(string(b)))
	if err != nil {
		log.Printf("slack post failed: %v", err)
	}
}

func main() {
	// Env (wired in docker-compose)
	endpoint := mustEnv("MINIO_ENDPOINT", "http://minio:9000")
	accessKey := mustEnv("MINIO_ACCESS_KEY", "minio")
	secretKey := mustEnv("MINIO_SECRET_KEY", "minio12345")
	bucket := mustEnv("MINIO_BUCKET", "mobility-delta")
	prefix := mustEnv("WATCH_PREFIX", "bronze/trips")
	webhook := strings.TrimSpace(os.Getenv("SLACK_WEBHOOK_URL"))
	slaStr := mustEnv("FRESHNESS_SLA_SECONDS", "900") // 15 min default
	slaSeconds, _ := strconv.Atoi(slaStr)

	host, secure := parseSecure(endpoint)
	client, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		log.Fatalf("minio connect: %v", err)
	}

	// Prometheus
	prometheus.MustRegister(metricObjCount, metricLastMod, metricFreshness)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("metrics listening on :9108/metrics")
		_ = http.ListenAndServe(":9108", nil)
	}()

	// Periodic scan
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	var lastAlert string
	for {
		now := time.Now().UTC()
		objCount := 0
		newest := time.Time{}
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)

		for obj := range client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		}) {
			if obj.Err != nil {
				log.Printf("list err: %v", obj.Err)
				continue
			}
			objCount++
			if obj.LastModified.After(newest) {
				newest = obj.LastModified
			}
		}
		cancel()

		metricObjCount.Set(float64(objCount))
		if !newest.IsZero() {
			metricLastMod.Set(float64(newest.Unix()))
			fresh := now.Sub(newest).Seconds()
			metricFreshness.Set(fresh)

			// Optional Slack alert on SLA breach (dedupe same minute msg)
			if webhook != "" && slaSeconds > 0 && int(fresh) > slaSeconds {
				msg := "⏰ Freshness breach: " + strconv.Itoa(int(fresh)) + "s > SLA " + strconv.Itoa(slaSeconds) +
					" (bucket=" + bucket + ", prefix=" + prefix + ")"
				if minute := now.Format("2006-01-02T15:04"); minute != lastAlert { // simple dedupe
					postSlack(webhook, msg)
					lastAlert = minute
				}
			}
		} else {
			metricLastMod.Set(0)
			metricFreshness.Set(1e12)
		}

		log.Printf("scan ok: count=%d newest=%v", objCount, newest)
		<-ticker.C
	}
}
