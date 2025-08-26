package cgroup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Scanner handles cgroup discovery and scanning
type Scanner struct {
	cgroupPath string
	logger     *logrus.Logger
	maxCgroups int
}

// CgroupInfo represents information about a cgroup
type CgroupInfo struct {
	Path        string
	Name        string
	Controllers []string
	Processes   []int
	LastScanned time.Time
}

// NewScanner creates a new cgroup scanner
func NewScanner(cgroupPath string, logger *logrus.Logger) *Scanner {
	return &Scanner{
		cgroupPath: cgroupPath,
		logger:     logger,
		maxCgroups: 10000,
	}
}

// Scan discovers and returns information about all cgroups
func (s *Scanner) Scan(ctx context.Context) ([]*CgroupInfo, error) {
	var cgroups []*CgroupInfo

	err := filepath.Walk(s.cgroupPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking on errors
		}

		if !info.IsDir() {
			return nil
		}

		// Skip if we've reached the maximum number of cgroups
		if len(cgroups) >= s.maxCgroups {
			return filepath.SkipDir
		}

		// Check if this directory contains cgroup.controllers file (indicates it's a cgroup)
		controllersFile := filepath.Join(path, "cgroup.controllers")
		if _, err := os.Stat(controllersFile); os.IsNotExist(err) {
			return nil
		}

		// Read controllers
		controllers, err := s.readControllers(controllersFile)
		if err != nil {
			s.logger.WithError(err).WithField("path", path).Debug("Failed to read controllers")
			return nil
		}

		// Create cgroup info
		cgroupInfo := &CgroupInfo{
			Path:        path,
			Name:        s.getCgroupName(path),
			Controllers: controllers,
			LastScanned: time.Now(),
		}

		cgroups = append(cgroups, cgroupInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan cgroups: %w", err)
	}

	s.logger.WithField("count", len(cgroups)).Debug("Scanned cgroups")
	return cgroups, nil
}

// readControllers reads the list of available controllers from cgroup.controllers file
func (s *Scanner) readControllers(controllersFile string) ([]string, error) {
	data, err := os.ReadFile(controllersFile)
	if err != nil {
		return nil, err
	}

	controllers := strings.Fields(strings.TrimSpace(string(data)))
	return controllers, nil
}

// getCgroupName extracts a readable name from the cgroup path
func (s *Scanner) getCgroupName(path string) string {
	// Remove the base cgroup path to get relative path
	relativePath := strings.TrimPrefix(path, s.cgroupPath)
	relativePath = strings.TrimPrefix(relativePath, "/")

	if relativePath == "" {
		return "root"
	}

	// Replace slashes with dots for metric labels
	return strings.ReplaceAll(relativePath, "/", ".")
}

// SetMaxCgroups sets the maximum number of cgroups to scan
func (s *Scanner) SetMaxCgroups(max int) {
	s.maxCgroups = max
}
