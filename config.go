package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the user preferences for the localize app.
type Config struct {
	Cities []string `json:"cities"` // List of selected city names
	Preset string   `json:"preset"` // Selected preset name
}

// DefaultConfigPath returns the default config file path (~/.localize/config.json).
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "config.json"
	}
	return filepath.Join(home, ".localize", "config.json")
}

// EnsureConfigDir creates the ~/.localize directory if it doesn't exist.
func EnsureConfigDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(home, ".localize")
	return os.MkdirAll(configDir, 0755)
}

// LoadConfig reads the configuration from the default config file.
// If the file doesn't exist, it returns a default config.
func LoadConfig() (*Config, error) {
	configPath := DefaultConfigPath()

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return &Config{}, nil
	}

	// Read and parse the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig writes the current configuration to the default config file.
func SaveConfig(config *Config) error {
	// Ensure the config directory exists
	if err := EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to config file
	configPath := DefaultConfigPath()
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// AvailableCity represents a city that can be displayed.
type AvailableCity struct {
	Name     string
	Timezone string
	Region   string // "americas", "europe", "asia", "africa", "oceania"
}

// GetAvailableCities returns a list of all supported cities with their timezones.
func GetAvailableCities() []AvailableCity {
	return []AvailableCity{
		// Americas
		{Name: "Honolulu", Timezone: "Pacific/Honolulu", Region: "americas"},
		{Name: "Anchorage", Timezone: "America/Anchorage", Region: "americas"},
		{Name: "Los Angeles", Timezone: "America/Los_Angeles", Region: "americas"},
		{Name: "New York", Timezone: "America/New_York", Region: "americas"},
		{Name: "Sao Paulo", Timezone: "America/Sao_Paulo", Region: "americas"},
		// Europe
		{Name: "London", Timezone: "Europe/London", Region: "europe"},
		{Name: "Paris", Timezone: "Europe/Paris", Region: "europe"},
		{Name: "Moscow", Timezone: "Europe/Moscow", Region: "europe"},
		// Africa
		{Name: "Cairo", Timezone: "Africa/Cairo", Region: "africa"},
		{Name: "Nairobi", Timezone: "Africa/Nairobi", Region: "africa"},
		// Asia
		{Name: "Dubai", Timezone: "Asia/Dubai", Region: "asia"},
		{Name: "Mumbai", Timezone: "Asia/Kolkata", Region: "asia"},
		{Name: "Singapore", Timezone: "Asia/Singapore", Region: "asia"},
		{Name: "Shanghai", Timezone: "Asia/Shanghai", Region: "asia"},
		{Name: "Tokyo", Timezone: "Asia/Tokyo", Region: "asia"},
		// Oceania
		{Name: "Sydney", Timezone: "Australia/Sydney", Region: "oceania"},
	}
}
