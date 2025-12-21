package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	GQLResolverDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "lighthouse",
			Subsystem: "graphql",
			Name:      "resolver_duration_seconds",
			Help:      "GraphQL resolver latency",
			// 默认 buckets：0.005s ~ 10s，足够用了
			Buckets: prometheus.DefBuckets,
		},
		[]string{"object", "field"},
	)

	GQLOperationTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "lighthouse",
			Subsystem: "graphql",
			Name:      "operations_total",
			Help:      "GraphQL operations count",
		},
		[]string{"operation", "type"},
	)
)

func Init() {
	prometheus.MustRegister(
		GQLResolverDuration,
		GQLOperationTotal,
	)
}
