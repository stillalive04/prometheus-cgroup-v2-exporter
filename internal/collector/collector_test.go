package collector

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stillalive04/prometheus-cgroup-v2-exporter/internal/config"
)

func TestNewCollectors(t *testing.T) {
	cfg := &config.Config{
		Collectors: config.CollectorsConfig{
			CPU: config.CPUCollectorConfig{
				Enabled:         true,
				IncludePressure: true,
			},
			Memory: config.MemoryCollectorConfig{
				Enabled:         true,
				IncludePressure: true,
				IncludeSwap:     true,
			},
			IO: config.IOCollectorConfig{
				Enabled:         true,
				IncludePressure: true,
				Devices:         []string{},
			},
			PIDs: config.PIDsCollectorConfig{
				Enabled: true,
			},
		},
		Advanced: config.AdvancedConfig{
			CacheDuration: 60,
		},
		Cgroup: config.CgroupConfig{
			Path: "/sys/fs/cgroup",
		},
	}

	logger := logrus.New()

	collectors, err := NewCollectors(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create collectors: %v", err)
	}

	if len(collectors) != 4 {
		t.Errorf("Expected 4 collectors, got %d", len(collectors))
	}

	expectedCollectors := []string{"cpu", "memory", "io", "pids"}
	for _, name := range expectedCollectors {
		if _, exists := collectors[name]; !exists {
			t.Errorf("Expected collector '%s' not found", name)
		}
	}
}

func TestNewCollectors_NoEnabled(t *testing.T) {
	cfg := &config.Config{
		Collectors: config.CollectorsConfig{
			CPU: config.CPUCollectorConfig{
				Enabled: false,
			},
			Memory: config.MemoryCollectorConfig{
				Enabled: false,
			},
			IO: config.IOCollectorConfig{
				Enabled: false,
			},
			PIDs: config.PIDsCollectorConfig{
				Enabled: false,
			},
		},
		Advanced: config.AdvancedConfig{
			CacheDuration: 60,
		},
		Cgroup: config.CgroupConfig{
			Path: "/sys/fs/cgroup",
		},
	}

	logger := logrus.New()

	_, err := NewCollectors(cfg, logger)
	if err == nil {
		t.Error("Expected error when no collectors are enabled")
	}
}

func TestBaseCollector(t *testing.T) {
	cfg := &config.Config{
		Advanced: config.AdvancedConfig{
			CacheDuration: 60,
		},
		Cgroup: config.CgroupConfig{
			Path: "/sys/fs/cgroup",
		},
	}

	logger := logrus.New()

	bc := NewBaseCollector("test", true, cfg, logger)

	if bc.Name() != "test" {
		t.Errorf("Expected name 'test', got '%s'", bc.Name())
	}

	if !bc.Enabled() {
		t.Error("Expected collector to be enabled")
	}

	// Test cache functionality
	key := "test-key"
	value := "test-value"

	// Initially, cache should be empty
	if _, exists := bc.GetCachedData(key); exists {
		t.Error("Cache should be empty initially")
	}

	// Set cached data
	bc.SetCachedData(key, value)

	// Retrieve cached data
	if cachedValue, exists := bc.GetCachedData(key); !exists {
		t.Error("Cached data should exist")
	} else if cachedValue != value {
		t.Errorf("Expected cached value '%s', got '%s'", value, cachedValue)
	}

	// Clear cache
	bc.ClearCache()

	// Cache should be empty again
	if _, exists := bc.GetCachedData(key); exists {
		t.Error("Cache should be empty after clearing")
	}
}
