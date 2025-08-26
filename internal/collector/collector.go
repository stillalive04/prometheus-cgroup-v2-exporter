package collector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/stillalive04/prometheus-cgroup-v2-exporter/internal/config"
	"github.com/stillalive04/prometheus-cgroup-v2-exporter/internal/cgroup"
)

// Collector interface defines the contract for all collectors
type Collector interface {
	prometheus.Collector
	Name() string
	Enabled() bool
}

// BaseCollector provides common functionality for all collectors
type BaseCollector struct {
	name    string
	enabled bool
	config  *config.Config
	logger  *logrus.Logger
	scanner *cgroup.Scanner
	mutex   sync.RWMutex
	
	// Metrics cache
	cache      map[string]interface{}
	cacheTime  time.Time
	cacheTTL   time.Duration
}

// NewBaseCollector creates a new base collector
func NewBaseCollector(name string, enabled bool, cfg *config.Config, logger *logrus.Logger) *BaseCollector {
	scanner := cgroup.NewScanner(cfg.Cgroup.Path, logger)
	
	return &BaseCollector{
		name:      name,
		enabled:   enabled,
		config:    cfg,
		logger:    logger,
		scanner:   scanner,
		cache:     make(map[string]interface{}),
		cacheTTL:  cfg.Advanced.CacheDuration,
	}
}

// Name returns the collector name
func (bc *BaseCollector) Name() string {
	return bc.name
}

// Enabled returns whether the collector is enabled
func (bc *BaseCollector) Enabled() bool {
	return bc.enabled
}

// GetCachedData retrieves cached data if still valid
func (bc *BaseCollector) GetCachedData(key string) (interface{}, bool) {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()
	
	if time.Since(bc.cacheTime) > bc.cacheTTL {
		return nil, false
	}
	
	data, exists := bc.cache[key]
	return data, exists
}

// SetCachedData stores data in cache
func (bc *BaseCollector) SetCachedData(key string, data interface{}) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	
	bc.cache[key] = data
	bc.cacheTime = time.Now()
}

// ClearCache clears the cache
func (bc *BaseCollector) ClearCache() {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	
	bc.cache = make(map[string]interface{})
}

// NewCollectors creates and returns all enabled collectors
func NewCollectors(cfg *config.Config, logger *logrus.Logger) (map[string]Collector, error) {
	collectors := make(map[string]Collector)
	
	// CPU Collector
	if cfg.Collectors.CPU.Enabled {
		cpuCollector := NewCPUCollector(cfg, logger)
		collectors["cpu"] = cpuCollector
	}
	
	// Memory Collector
	if cfg.Collectors.Memory.Enabled {
		memoryCollector := NewMemoryCollector(cfg, logger)
		collectors["memory"] = memoryCollector
	}
	
	// I/O Collector
	if cfg.Collectors.IO.Enabled {
		ioCollector := NewIOCollector(cfg, logger)
		collectors["io"] = ioCollector
	}
	
	// PIDs Collector
	if cfg.Collectors.PIDs.Enabled {
		pidsCollector := NewPIDsCollector(cfg, logger)
		collectors["pids"] = pidsCollector
	}
	
	if len(collectors) == 0 {
		return nil, fmt.Errorf("no collectors enabled")
	}
	
	logger.WithField("collectors", len(collectors)).Info("Initialized collectors")
	return collectors, nil
}

// CollectorMetrics holds common metrics for all collectors
type CollectorMetrics struct {
	ScrapeDuration prometheus.Histogram
	ScrapeErrors   prometheus.Counter
	LastScrapeTime prometheus.Gauge
	CgroupsScraped prometheus.Gauge
}

// NewCollectorMetrics creates common metrics for a collector
func NewCollectorMetrics(subsystem string) *CollectorMetrics {
	return &CollectorMetrics{
		ScrapeDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "prometheus_cgroup_v2_exporter",
			Subsystem: subsystem,
			Name:      "scrape_duration_seconds",
			Help:      "Time spent scraping metrics",
			Buckets:   prometheus.DefBuckets,
		}),
		ScrapeErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "prometheus_cgroup_v2_exporter",
			Subsystem: subsystem,
			Name:      "scrape_errors_total",
			Help:      "Total number of scrape errors",
		}),
		LastScrapeTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "prometheus_cgroup_v2_exporter",
			Subsystem: subsystem,
			Name:      "last_scrape_timestamp_seconds",
			Help:      "Unix timestamp of the last successful scrape",
		}),
		CgroupsScraped: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "prometheus_cgroup_v2_exporter",
			Subsystem: subsystem,
			Name:      "cgroups_scraped",
			Help:      "Number of cgroups scraped in the last collection",
		}),
	}
}

// Describe implements prometheus.Collector
func (cm *CollectorMetrics) Describe(ch chan<- *prometheus.Desc) {
	cm.ScrapeDuration.Describe(ch)
	cm.ScrapeErrors.Describe(ch)
	cm.LastScrapeTime.Describe(ch)
	cm.CgroupsScraped.Describe(ch)
}

// Collect implements prometheus.Collector
func (cm *CollectorMetrics) Collect(ch chan<- prometheus.Metric) {
	cm.ScrapeDuration.Collect(ch)
	cm.ScrapeErrors.Collect(ch)
	cm.LastScrapeTime.Collect(ch)
	cm.CgroupsScraped.Collect(ch)
}
