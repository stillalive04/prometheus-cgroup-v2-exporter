package config

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Test loading default configuration
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load default config: %v", err)
	}

	// Verify default values
	if cfg.Web.ListenAddress != ":9753" {
		t.Errorf("Expected default listen address ':9753', got '%s'", cfg.Web.ListenAddress)
	}

	if cfg.Web.TelemetryPath != "/metrics" {
		t.Errorf("Expected default telemetry path '/metrics', got '%s'", cfg.Web.TelemetryPath)
	}

	if cfg.Cgroup.Path != "/sys/fs/cgroup" {
		t.Errorf("Expected default cgroup path '/sys/fs/cgroup', got '%s'", cfg.Cgroup.Path)
	}

	if !cfg.Collectors.CPU.Enabled {
		t.Error("CPU collector should be enabled by default")
	}

	if !cfg.Collectors.Memory.Enabled {
		t.Error("Memory collector should be enabled by default")
	}

	if !cfg.Collectors.IO.Enabled {
		t.Error("I/O collector should be enabled by default")
	}

	if !cfg.Collectors.PIDs.Enabled {
		t.Error("PIDs collector should be enabled by default")
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("Expected default log level 'info', got '%s'", cfg.Logging.Level)
	}

	if cfg.Advanced.MaxCgroups != 10000 {
		t.Errorf("Expected default max cgroups 10000, got %d", cfg.Advanced.MaxCgroups)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Web: WebConfig{
					ListenAddress: ":9753",
					TelemetryPath: "/metrics",
				},
				Cgroup: CgroupConfig{
					Path:            "/sys/fs/cgroup",
					RefreshInterval: 15 * time.Second,
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "logfmt",
				},
				Advanced: AdvancedConfig{
					MaxCgroups:    10000,
					ScanInterval:  30 * time.Second,
					CacheDuration: 60 * time.Second,
				},
			},
			wantErr: false,
		},
		{
			name: "empty listen address",
			config: &Config{
				Web: WebConfig{
					ListenAddress: "",
					TelemetryPath: "/metrics",
				},
				Cgroup: CgroupConfig{
					Path:            "/sys/fs/cgroup",
					RefreshInterval: 15 * time.Second,
				},
				Logging: LoggingConfig{
					Level:  "info",
					Format: "logfmt",
				},
				Advanced: AdvancedConfig{
					MaxCgroups:    10000,
					ScanInterval:  30 * time.Second,
					CacheDuration: 60 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			config: &Config{
				Web: WebConfig{
					ListenAddress: ":9753",
					TelemetryPath: "/metrics",
				},
				Cgroup: CgroupConfig{
					Path:            "/sys/fs/cgroup",
					RefreshInterval: 15 * time.Second,
				},
				Logging: LoggingConfig{
					Level:  "invalid",
					Format: "logfmt",
				},
				Advanced: AdvancedConfig{
					MaxCgroups:    10000,
					ScanInterval:  30 * time.Second,
					CacheDuration: 60 * time.Second,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
