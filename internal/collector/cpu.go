package collector

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/stillalive04/prometheus-cgroup-v2-exporter/internal/config"
)

// CPUCollector collects CPU-related metrics from cgroup v2
type CPUCollector struct {
	*BaseCollector
	metrics *CollectorMetrics
	
	// CPU metrics
	cpuUsageTotal     *prometheus.CounterVec
	cpuUserTotal      *prometheus.CounterVec
	cpuSystemTotal    *prometheus.CounterVec
	cpuThrottledTotal *prometheus.CounterVec
	cpuPeriodsTotal   *prometheus.CounterVec
	cpuPressureTotal  *prometheus.CounterVec
}

// NewCPUCollector creates a new CPU collector
func NewCPUCollector(cfg *config.Config, logger *logrus.Logger) *CPUCollector {
	base := NewBaseCollector("cpu", cfg.Collectors.CPU.Enabled, cfg, logger)
	
	collector := &CPUCollector{
		BaseCollector: base,
		metrics:       NewCollectorMetrics("cpu"),
	}
	
	// Initialize metrics
	collector.initMetrics()
	
	return collector
}

func (c *CPUCollector) initMetrics() {
	c.cpuUsageTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cgroup",
			Subsystem: "cpu",
			Name:      "usage_seconds_total",
			Help:      "Total CPU time consumed by cgroup",
		},
		[]string{"cgroup", "mode"},
	)
	
	c.cpuUserTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cgroup",
			Subsystem: "cpu",
			Name:      "user_seconds_total",
			Help:      "Total CPU time spent in user mode by cgroup",
		},
		[]string{"cgroup"},
	)
	
	c.cpuSystemTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cgroup",
			Subsystem: "cpu",
			Name:      "system_seconds_total",
			Help:      "Total CPU time spent in system mode by cgroup",
		},
		[]string{"cgroup"},
	)
	
	c.cpuThrottledTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cgroup",
			Subsystem: "cpu",
			Name:      "throttled_seconds_total",
			Help:      "Total time spent throttled by cgroup",
		},
		[]string{"cgroup"},
	)
	
	c.cpuPeriodsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cgroup",
			Subsystem: "cpu",
			Name:      "periods_total",
			Help:      "Total number of CPU periods by cgroup",
		},
		[]string{"cgroup"},
	)
	
	if c.config.Collectors.CPU.IncludePressure {
		c.cpuPressureTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "cgroup",
				Subsystem: "cpu",
				Name:      "pressure_seconds_total",
				Help:      "Total CPU pressure stall time by cgroup",
			},
			[]string{"cgroup", "type"},
		)
	}
}

// Describe implements prometheus.Collector
func (c *CPUCollector) Describe(ch chan<- *prometheus.Desc) {
	if !c.Enabled() {
		return
	}
	
	c.cpuUsageTotal.Describe(ch)
	c.cpuUserTotal.Describe(ch)
	c.cpuSystemTotal.Describe(ch)
	c.cpuThrottledTotal.Describe(ch)
	c.cpuPeriodsTotal.Describe(ch)
	
	if c.cpuPressureTotal != nil {
		c.cpuPressureTotal.Describe(ch)
	}
	
	c.metrics.Describe(ch)
}

// Collect implements prometheus.Collector
func (c *CPUCollector) Collect(ch chan<- prometheus.Metric) {
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
	c.cpuUsageTotal.Collect(ch)
	c.cpuUserTotal.Collect(ch)
	c.cpuSystemTotal.Collect(ch)
	c.cpuThrottledTotal.Collect(ch)
	c.cpuPeriodsTotal.Collect(ch)
	
	if c.cpuPressureTotal != nil {
		c.cpuPressureTotal.Collect(ch)
	}
	
	c.metrics.Collect(ch)
}

func (c *CPUCollector) collectCgroupMetrics(cgroup interface{}) {
	// This is a stub implementation
	// In a real implementation, this would read from cgroup v2 files
	// like cpu.stat, cpu.pressure, etc.
	
	// For now, just set some dummy metrics to make the code compile
	c.cpuUsageTotal.WithLabelValues("example", "user").Add(0)
	c.cpuUserTotal.WithLabelValues("example").Add(0)
	c.cpuSystemTotal.WithLabelValues("example").Add(0)
}
