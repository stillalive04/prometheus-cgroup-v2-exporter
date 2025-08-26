package collector

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/stillalive04/prometheus-cgroup-v2-exporter/internal/config"
)

// IOCollector collects I/O-related metrics from cgroup v2
type IOCollector struct {
	*BaseCollector
	metrics *CollectorMetrics
	
	// I/O metrics
	ioReadBytesTotal  *prometheus.CounterVec
	ioWriteBytesTotal *prometheus.CounterVec
	ioReadOpsTotal    *prometheus.CounterVec
	ioWriteOpsTotal   *prometheus.CounterVec
	ioPressureTotal   *prometheus.CounterVec
}

// NewIOCollector creates a new I/O collector
func NewIOCollector(cfg *config.Config, logger *logrus.Logger) *IOCollector {
	base := NewBaseCollector("io", cfg.Collectors.IO.Enabled, cfg, logger)
	
	collector := &IOCollector{
		BaseCollector: base,
		metrics:       NewCollectorMetrics("io"),
	}
	
	// Initialize metrics
	collector.initMetrics()
	
	return collector
}

func (c *IOCollector) initMetrics() {
	c.ioReadBytesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cgroup",
			Subsystem: "io",
			Name:      "read_bytes_total",
			Help:      "Total bytes read by cgroup",
		},
		[]string{"cgroup", "device"},
	)
	
	c.ioWriteBytesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cgroup",
			Subsystem: "io",
			Name:      "write_bytes_total",
			Help:      "Total bytes written by cgroup",
		},
		[]string{"cgroup", "device"},
	)
	
	c.ioReadOpsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cgroup",
			Subsystem: "io",
			Name:      "read_operations_total",
			Help:      "Total read operations by cgroup",
		},
		[]string{"cgroup", "device"},
	)
	
	c.ioWriteOpsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cgroup",
			Subsystem: "io",
			Name:      "write_operations_total",
			Help:      "Total write operations by cgroup",
		},
		[]string{"cgroup", "device"},
	)
	
	if c.config.Collectors.IO.IncludePressure {
		c.ioPressureTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "cgroup",
				Subsystem: "io",
				Name:      "pressure_seconds_total",
				Help:      "Total I/O pressure stall time by cgroup",
			},
			[]string{"cgroup", "type"},
		)
	}
}

// Describe implements prometheus.Collector
func (c *IOCollector) Describe(ch chan<- *prometheus.Desc) {
	if !c.Enabled() {
		return
	}
	
	c.ioReadBytesTotal.Describe(ch)
	c.ioWriteBytesTotal.Describe(ch)
	c.ioReadOpsTotal.Describe(ch)
	c.ioWriteOpsTotal.Describe(ch)
	
	if c.ioPressureTotal != nil {
		c.ioPressureTotal.Describe(ch)
	}
	
	c.metrics.Describe(ch)
}

// Collect implements prometheus.Collector
func (c *IOCollector) Collect(ch chan<- prometheus.Metric) {
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
	c.ioReadBytesTotal.Collect(ch)
	c.ioWriteBytesTotal.Collect(ch)
	c.ioReadOpsTotal.Collect(ch)
	c.ioWriteOpsTotal.Collect(ch)
	
	if c.ioPressureTotal != nil {
		c.ioPressureTotal.Collect(ch)
	}
	
	c.metrics.Collect(ch)
}

func (c *IOCollector) collectCgroupMetrics(cgroup interface{}) {
	// This is a stub implementation
	// In a real implementation, this would read from cgroup v2 files
	// like io.stat, io.pressure, etc.
	
	// For now, just set some dummy metrics to make the code compile
	c.ioReadBytesTotal.WithLabelValues("example", "sda").Add(0)
	c.ioWriteBytesTotal.WithLabelValues("example", "sda").Add(0)
}
