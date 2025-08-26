package collector

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/stillalive04/prometheus-cgroup-v2-exporter/internal/config"
)

// PIDsCollector collects process-related metrics from cgroup v2
type PIDsCollector struct {
	*BaseCollector
	metrics *CollectorMetrics

	// PIDs metrics
	processesCount   *prometheus.GaugeVec
	processesRunning *prometheus.GaugeVec
	processesSleeping *prometheus.GaugeVec
	processesZombie  *prometheus.GaugeVec
}

// NewPIDsCollector creates a new PIDs collector
func NewPIDsCollector(cfg *config.Config, logger *logrus.Logger) *PIDsCollector {
	base := NewBaseCollector("pids", cfg.Collectors.PIDs.Enabled, cfg, logger)

	collector := &PIDsCollector{
		BaseCollector: base,
		metrics:       NewCollectorMetrics("pids"),
	}

	// Initialize metrics
	collector.initMetrics()

	return collector
}

func (c *PIDsCollector) initMetrics() {
	c.processesCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup",
			Subsystem: "processes",
			Name:      "count",
			Help:      "Number of processes in cgroup",
		},
		[]string{"cgroup"},
	)
	
	c.processesRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup",
			Subsystem: "processes",
			Name:      "running",
			Help:      "Number of running processes in cgroup",
		},
		[]string{"cgroup"},
	)
	
	c.processesSleeping = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup",
			Subsystem: "processes",
			Name:      "sleeping",
			Help:      "Number of sleeping processes in cgroup",
		},
		[]string{"cgroup"},
	)
	
	c.processesZombie = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "cgroup",
			Subsystem: "processes",
			Name:      "zombie",
			Help:      "Number of zombie processes in cgroup",
		},
		[]string{"cgroup"},
	)
}

// Describe implements prometheus.Collector
func (c *PIDsCollector) Describe(ch chan<- *prometheus.Desc) {
	if !c.Enabled() {
		return
	}
	
	c.processesCount.Describe(ch)
	c.processesRunning.Describe(ch)
	c.processesSleeping.Describe(ch)
	c.processesZombie.Describe(ch)
	
	c.metrics.Describe(ch)
}

// Collect implements prometheus.Collector
func (c *PIDsCollector) Collect(ch chan<- prometheus.Metric) {
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
	c.processesCount.Collect(ch)
	c.processesRunning.Collect(ch)
	c.processesSleeping.Collect(ch)
	c.processesZombie.Collect(ch)
	
	c.metrics.Collect(ch)
}

func (c *PIDsCollector) collectCgroupMetrics(cgroup interface{}) {
	// This is a stub implementation
	// In a real implementation, this would read from cgroup v2 files
	// like cgroup.procs, pids.current, etc.
	
	// For now, just set some dummy metrics to make the code compile
	c.processesCount.WithLabelValues("example").Set(0)
	c.processesRunning.WithLabelValues("example").Set(0)
}
