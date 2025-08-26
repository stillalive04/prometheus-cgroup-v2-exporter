package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Web        WebConfig        `mapstructure:"web"`
	Cgroup     CgroupConfig     `mapstructure:"cgroup"`
	Collectors CollectorsConfig `mapstructure:"collectors"`
	Logging    LoggingConfig    `mapstructure:"logging"`
	Advanced   AdvancedConfig   `mapstructure:"advanced"`
}

// WebConfig contains web server configuration
type WebConfig struct {
	ListenAddress string `mapstructure:"listen_address"`
	TelemetryPath string `mapstructure:"telemetry_path"`
}

// CgroupConfig contains cgroup-related configuration
type CgroupConfig struct {
	Path            string        `mapstructure:"path"`
	RefreshInterval time.Duration `mapstructure:"refresh_interval"`
}

// CollectorsConfig contains collector configuration
type CollectorsConfig struct {
	CPU    CPUCollectorConfig    `mapstructure:"cpu"`
	Memory MemoryCollectorConfig `mapstructure:"memory"`
	IO     IOCollectorConfig     `mapstructure:"io"`
	PIDs   PIDsCollectorConfig   `mapstructure:"pids"`
}

// CPUCollectorConfig contains CPU collector configuration
type CPUCollectorConfig struct {
	Enabled         bool `mapstructure:"enabled"`
	IncludePressure bool `mapstructure:"include_pressure"`
}

// MemoryCollectorConfig contains memory collector configuration
type MemoryCollectorConfig struct {
	Enabled         bool `mapstructure:"enabled"`
	IncludePressure bool `mapstructure:"include_pressure"`
	IncludeSwap     bool `mapstructure:"include_swap"`
}

// IOCollectorConfig contains I/O collector configuration
type IOCollectorConfig struct {
	Enabled         bool     `mapstructure:"enabled"`
	IncludePressure bool     `mapstructure:"include_pressure"`
	Devices         []string `mapstructure:"devices"`
}

// PIDsCollectorConfig contains PIDs collector configuration
type PIDsCollectorConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// AdvancedConfig contains advanced configuration options
type AdvancedConfig struct {
	MaxCgroups    int           `mapstructure:"max_cgroups"`
	ScanInterval  time.Duration `mapstructure:"scan_interval"`
	CacheDuration time.Duration `mapstructure:"cache_duration"`
}

// Load loads configuration from various sources
func Load() (*Config, error) {
	// Set defaults
	setDefaults()

	// Read configuration file if specified
	if configFile := viper.GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal configuration
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// Web defaults
	viper.SetDefault("web.listen_address", ":9753")
	viper.SetDefault("web.telemetry_path", "/metrics")

	// Cgroup defaults
	viper.SetDefault("cgroup.path", "/sys/fs/cgroup")
	viper.SetDefault("cgroup.refresh_interval", "15s")

	// Collector defaults
	viper.SetDefault("collectors.cpu.enabled", true)
	viper.SetDefault("collectors.cpu.include_pressure", true)
	viper.SetDefault("collectors.memory.enabled", true)
	viper.SetDefault("collectors.memory.include_pressure", true)
	viper.SetDefault("collectors.memory.include_swap", true)
	viper.SetDefault("collectors.io.enabled", true)
	viper.SetDefault("collectors.io.include_pressure", true)
	viper.SetDefault("collectors.io.devices", []string{})
	viper.SetDefault("collectors.pids.enabled", true)

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "logfmt")

	// Advanced defaults
	viper.SetDefault("advanced.max_cgroups", 10000)
	viper.SetDefault("advanced.scan_interval", "30s")
	viper.SetDefault("advanced.cache_duration", "60s")
}

func validate(config *Config) error {
	// Validate web configuration
	if config.Web.ListenAddress == "" {
		return fmt.Errorf("web.listen_address cannot be empty")
	}
	if config.Web.TelemetryPath == "" {
		return fmt.Errorf("web.telemetry_path cannot be empty")
	}

	// Validate cgroup configuration
	if config.Cgroup.Path == "" {
		return fmt.Errorf("cgroup.path cannot be empty")
	}
	if config.Cgroup.RefreshInterval <= 0 {
		return fmt.Errorf("cgroup.refresh_interval must be positive")
	}

	// Validate logging configuration
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[config.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", config.Logging.Level)
	}

	validLogFormats := map[string]bool{
		"logfmt": true, "json": true,
	}
	if !validLogFormats[config.Logging.Format] {
		return fmt.Errorf("invalid log format: %s", config.Logging.Format)
	}

	// Validate advanced configuration
	if config.Advanced.MaxCgroups <= 0 {
		return fmt.Errorf("advanced.max_cgroups must be positive")
	}
	if config.Advanced.ScanInterval <= 0 {
		return fmt.Errorf("advanced.scan_interval must be positive")
	}
	if config.Advanced.CacheDuration <= 0 {
		return fmt.Errorf("advanced.cache_duration must be positive")
	}

	return nil
}
