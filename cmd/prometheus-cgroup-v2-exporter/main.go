package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stillalive04/prometheus-cgroup-v2-exporter/internal/collector"
	"github.com/stillalive04/prometheus-cgroup-v2-exporter/internal/config"
)

var (
	cfg     *config.Config
	log     = logrus.New()
	rootCmd = &cobra.Command{
		Use:   "prometheus-cgroup-v2-exporter",
		Short: "Prometheus exporter for cgroup v2 metrics",
		Long: `A high-performance Prometheus exporter for cgroup v2 metrics,
designed for modern containerized environments and Kubernetes clusters.`,
		Version: version.Version,
		RunE:    run,
	}
)

func init() {
	// Command line flags
	rootCmd.PersistentFlags().String("web.listen-address", ":9753", "Address to listen on for web interface and telemetry")
	rootCmd.PersistentFlags().String("web.telemetry-path", "/metrics", "Path under which to expose metrics")
	rootCmd.PersistentFlags().String("cgroup.path", "/sys/fs/cgroup", "Path to cgroup v2 filesystem")
	rootCmd.PersistentFlags().StringSlice("collector.enable", []string{"cpu", "memory", "io", "pids"}, "Comma-separated list of enabled collectors")
	rootCmd.PersistentFlags().StringSlice("collector.disable", []string{}, "Comma-separated list of disabled collectors")
	rootCmd.PersistentFlags().String("log.level", "info", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("log.format", "logfmt", "Log format (logfmt, json)")
	rootCmd.PersistentFlags().String("config", "", "Configuration file path")
	rootCmd.PersistentFlags().Duration("scan.interval", 30*time.Second, "Interval for scanning cgroups")
	rootCmd.PersistentFlags().Int("max.cgroups", 10000, "Maximum number of cgroups to monitor")
	rootCmd.PersistentFlags().Duration("cache.duration", 60*time.Second, "Duration to cache metrics")

	// Bind flags to viper
	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.SetEnvPrefix("CGROUPV2_EXPORTER")
	viper.AutomaticEnv()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Initialize configuration
	var err error
	cfg, err = config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup logging
	setupLogging()

	log.WithFields(logrus.Fields{
		"version":    version.Version,
		"revision":   version.Revision,
		"branch":     version.Branch,
		"build_date": version.BuildDate,
		"go_version": version.GoVersion,
	}).Info("Starting prometheus-cgroup-v2-exporter")

	// Validate cgroup v2 availability
	if err := validateCgroupV2(); err != nil {
		return fmt.Errorf("cgroup v2 validation failed: %w", err)
	}

	// Create Prometheus registry
	registry := prometheus.NewRegistry()

	// Register version metrics
	registry.MustRegister(version.NewCollector("prometheus_cgroup_v2_exporter"))

	// Initialize collectors
	collectors, err := collector.NewCollectors(cfg, log)
	if err != nil {
		return fmt.Errorf("failed to initialize collectors: %w", err)
	}

	// Register collectors
	for name, coll := range collectors {
		log.WithField("collector", name).Info("Registering collector")
		registry.MustRegister(coll)
	}

	// Setup HTTP server
	mux := http.NewServeMux()
	mux.Handle(cfg.Web.TelemetryPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:      log,
		ErrorHandling: promhttp.ContinueOnError,
	}))

	// Health check endpoint
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/ready", readyHandler)

	// Root handler with basic info
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<html>
<head><title>Prometheus cgroup v2 Exporter</title></head>
<body>
<h1>Prometheus cgroup v2 Exporter</h1>
<p><a href="%s">Metrics</a></p>
<p><a href="/health">Health</a></p>
<p><a href="/ready">Ready</a></p>
<p>Version: %s</p>
</body>
</html>`, cfg.Web.TelemetryPath, version.Version)
	})

	server := &http.Server{
		Addr:         cfg.Web.ListenAddress,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Info("Received shutdown signal")
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.WithError(err).Error("Server shutdown failed")
		}
	}()

	// Start collectors
	for name, coll := range collectors {
		if starter, ok := coll.(interface{ Start(context.Context) error }); ok {
			go func(name string, starter interface{ Start(context.Context) error }) {
				if err := starter.Start(ctx); err != nil {
					log.WithField("collector", name).WithError(err).Error("Collector failed")
				}
			}(name, starter)
		}
	}

	// Start HTTP server
	log.WithField("address", cfg.Web.ListenAddress).Info("Starting HTTP server")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server failed: %w", err)
	}

	log.Info("Exporter stopped")
	return nil
}

func setupLogging() {
	// Set log level
	level, err := logrus.ParseLevel(cfg.Logging.Level)
	if err != nil {
		log.WithError(err).Warn("Invalid log level, using info")
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	// Set log format
	if cfg.Logging.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}
}

func validateCgroupV2() error {
	// Check if cgroup v2 filesystem is mounted
	if _, err := os.Stat(cfg.Cgroup.Path); os.IsNotExist(err) {
		return fmt.Errorf("cgroup v2 filesystem not found at %s", cfg.Cgroup.Path)
	}

	// Check if it's actually cgroup v2
	cgroupVersion := cfg.Cgroup.Path + "/cgroup.controllers"
	if _, err := os.Stat(cgroupVersion); os.IsNotExist(err) {
		return fmt.Errorf("cgroup v2 controllers file not found, ensure cgroup v2 is enabled")
	}

	log.WithField("path", cfg.Cgroup.Path).Info("cgroup v2 filesystem validated")
	return nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	// Check if cgroup v2 is still accessible
	if err := validateCgroupV2(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, `{"status":"not ready","error":"%s","timestamp":"%s"}`, err.Error(), time.Now().Format(time.RFC3339))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ready","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}
