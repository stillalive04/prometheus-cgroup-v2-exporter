package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	// Basic test to ensure the package compiles
	if rootCmd == nil {
		t.Error("rootCmd should not be nil")
	}
}

func TestVersion(t *testing.T) {
	// Test version command
	rootCmd.SetArgs([]string{"--version"})

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Version command panicked: %v", r)
		}
	}()

	// We expect this to exit, so we don't actually execute it in tests
	// Just verify the command structure is valid
	if rootCmd.Use != "prometheus-cgroup-v2-exporter" {
		t.Errorf("Expected command name 'prometheus-cgroup-v2-exporter', got '%s'", rootCmd.Use)
	}
}
