package collector

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/stillalive04/prometheus-cgroup-v2-exporter/internal/config"
)

// MemoryCollector collects memory-related metrics from cgroup v2
type MemoryCollector struct {
	*BaseCollector
	metrics *CollectorMetrics

	// Memory metrics
	memoryUsageBytes     *prometheus.GaugeVec
	memoryLimitBytes     *prometheus.GaugeVec
	memoryCacheBytes     *prometheus.GaugeVec
	memoryRSSBytes       *prometheus.GaugeVec
	memorySwapUsageBytes *prometheus.GaugeVec
	memoryOOMEvents      *prometheus.CounterVec
	memoryPressureTotal  *prometheus.CounterVec
}

// NewMemoryCollector creates a new memory collector
func NewMemoryCollector(cfg *config.Config, logger *logrus.Logger) *MemoryCollector {
	base := NewBaseCollector("memory", cfg.Collectors.Memory.Enabled, cfg, logger)

	collector := &MemoryCollector{
		BaseCollector: base,
		metrics:       NewCollectorMetrics("memory"),
	}

	// Initialize metrics
	collector.initMetrics()

	return collector
}

func (c *MemoryCollector) initMetrics() {
	c.memoryUsageBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup",
			Subsystem: "memory",
			Name:      "usage_bytes",
			Help:      "Current memory usage by cgroup",
		},
		[]string{"cgroup"},
	)
	
	c.memoryLimitBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup",
			Subsystem: "memory",
			Name:      "limit_bytes",
			Help:      "Memory limit for cgroup",
		},
		[]string{"cgroup"},
	)
	
	c.memoryCacheBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup",
			Subsystem: "memory",
			Name:      "cache_bytes",
			Help:      "Cache memory usage by cgroup",
		},
		[]string{"cgroup"},
	)
	
	c.memoryRSSBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup",
			Subsystem: "memory",
			Name:      "rss_bytes",
			Help:      "RSS memory usage by cgroup",
		},
		[]string{"cgroup"},
	)
	
	if c.config.Collectors.Memory.IncludeSwap {
		c.memorySwapUsageBytes = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "cgroup",
				Subsystem: "memory",
				Name:      "swap_usage_bytes",
				Help:      "Swap usage by cgroup",
			},
			[]string{"cgroup"},
		)
	}
	
	c.memoryOOMEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cgroup",
			Subsystem: "memory",
			Name:      "oom_events_total",
			Help:      "Total number of OOM events by cgroup",
		},
		[]string{"cgroup"},
	)
	
	if c.config.Collectors.Memory.IncludePressure {
		c.memoryPressureTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "cgroup",
				Subsystem: "memory",
				Name:      "pressure_seconds_total",
				Help:      "Total memory pressure stall time by cgroup",
			},
			[]string{"cgroup", "type"},
		)
	}
}

// Describe implements prometheus.Collector
func (c *MemoryCollector) Describe(ch chan<- *prometheus.Desc) {
	if !c.Enabled() {
		return
	}
	
	c.memoryUsageBytes.Describe(ch)
	c.memoryLimitBytes.Describe(ch)
	c.memoryCacheBytes.Describe(ch)
	c.memoryRSSBytes.Describe(ch)
	c.memoryOOMEvents.Describe(ch)
	
	if c.memorySwapUsageBytes != nil {
		c.memorySwapUsageBytes.Describe(ch)
	}
	
	if c.memoryPressureTotal != nil {
		c.memoryPressureTotal.Describe(ch)
	}
	
	c.metrics.Describe(ch)
}

// Collect implements prometheus.Collector
func (c *MemoryCollector) Collect(ch chan<- prometheus.Metric) {
	if !c.Enabled() {
		return
	}
	
	start := time.Now()
	defer func() {
		c.metrics.ScrapeDuration.Observe(time.Since(start).Seconds())
		c.metrics.LastScrapeTime.SetToCurrentTime()
	}()
	
	// Scan cgroups
	cgroups, err := c.scanner.Scan(context.Background())
	if err != nil {
		c.logger.WithError(err).Error("Failed to scan cgroups")
		c.metrics.ScrapeErrors.Inc()
		return
	}
	
	c.metrics.CgroupsScraped.Set(float64(len(cgroups)))
	
	// Collect metrics from each cgroup
	for _, cgroup := range cgroups {
		c.collectCgroupMetrics(cgroup)
	}
	
	// Collect all metrics
	c.memoryUsageBytes.Collect(ch)
	c.memoryLimitBytes.Collect(ch)
	c.memoryCacheBytes.Collect(ch)
	c.memoryRSSBytes.Collect(ch)
	c.memoryOOMEvents.Collect(ch)
	
	if c.memorySwapUsageBytes != nil {
		c.memorySwapUsageBytes.Collect(ch)
	}
	
	if c.memoryPressureTotal != nil {
		c.memoryPressureTotal.Collect(ch)
	}
	
	c.metrics.Collect(ch)
}

func (c *MemoryCollector) collectCgroupMetrics(cgroup interface{}) {
	// This is a stub implementation
	// In a real implementation, this would read from cgroup v2 files
	// like memory.current, memory.max, memory.stat, etc.
	
	// For now, just set some dummy metrics to make the code compile
	c.memoryUsageBytes.WithLabelValues("example").Set(0)
	c.memoryLimitBytes.WithLabelValues("example").Set(0)
	c.memoryCacheBytes.WithLabelValues("example").Set(0)
}
