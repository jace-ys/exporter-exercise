package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	reServerVersion   = regexp.MustCompile(`redis_version:([0-9\.]+)`)
	reReplicationRole = regexp.MustCompile(`role:([a-z]+)`)
)

var (
	port      = kingpin.Flag("port", "Port binding for the exporter metrics server.").Envar("PORT").Default("9091").Int()
	namespace = kingpin.Flag("namespace", "Namespace prefix for exported metrics.").Envar("NAMESPACE").Default("").String()
	redisURL  = kingpin.Flag("redis.url", "Connection URL of the Redis server to collect metrics from.").Envar("REDIS_URL").Default("redis://127.0.0.1:6379").String()
)

func main() {
	kingpin.Parse()

	collector, err := NewRedisCollector(*namespace, *redisURL)
	if err != nil {
		log.Fatalf("error initialising collector: %s", err)
	}

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
		collector,
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	log.Printf("exporter listening on :%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

var _ prometheus.Collector = (*RedisCollector)(nil)

type RedisCollector struct {
	metricDescKeys             *prometheus.Desc
	metricDescBuildVersionInfo *prometheus.Desc
	redis                      *redis.Client
}

func NewRedisCollector(namespace, redisURL string) (*RedisCollector, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	return &RedisCollector{
		metricDescBuildVersionInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "redis", "build_version_info"),
			"Info of the current version of Redis running",
			[]string{"version", "role"}, nil,
		),
		metricDescKeys: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "redis", "keys"),
			"Gauge of the total count of keys at a given time",
			nil, nil,
		),
		redis: redis.NewClient(opts),
	}, nil
}

func (c *RedisCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *RedisCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	info, err := c.fetchBuildVersionInfo(ctx)
	if err != nil {
		log.Printf("error fetching build version info: %s\n", err)
	} else {
		log.Printf("fetched build version info: version=%s, role=%s", info.version, info.role)
		ch <- prometheus.MustNewConstMetric(c.metricDescBuildVersionInfo, prometheus.GaugeValue, 1, info.version, info.role)
	}

	keys, err := c.fetchKeyCount(ctx)
	if err != nil {
		log.Printf("error fetching key count: %s\n", err)
	} else {
		log.Printf("fetched key count: %d", keys)
		ch <- prometheus.MustNewConstMetric(c.metricDescKeys, prometheus.GaugeValue, float64(keys))
	}
}

type buildVersionInfo struct {
	version string
	role    string
}

func (c *RedisCollector) fetchBuildVersionInfo(ctx context.Context) (*buildVersionInfo, error) {
	info, err := c.redis.Info(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("error fetching info: %s", err)
	}

	versionResult := reServerVersion.FindStringSubmatch(info)
	if len(versionResult) != 2 {
		return nil, fmt.Errorf("error parsing server version from info")
	}

	roleResult := reReplicationRole.FindStringSubmatch(info)
	if len(roleResult) != 2 {
		return nil, fmt.Errorf("error parsing role from info")
	}

	return &buildVersionInfo{
		version: versionResult[1],
		role:    roleResult[1],
	}, nil
}

func (c *RedisCollector) fetchKeyCount(ctx context.Context) (int, error) {
	keys, err := c.redis.Keys(ctx, "*").Result()
	if err != nil {
		return 0, fmt.Errorf("error fetching keys: %s", err)
	}

	return len(keys), nil
}
