package config

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
	"task-manager/pkg/db"
	"task-manager/pkg/logger"
)

// Config represents the main application configuration.
// It can be loaded from a YAML file and/or overridden by environment variables.
type Config struct {
	Logger       logger.Config   `json:"logger" yaml:"LOGGER"`                 // Logger configuration
	DB           db.Configs      `json:"db" yaml:"DB"`                         // Database connection settings
	Redis        RedisConfig     `json:"redis" yaml:"REDIS"`                   // Redis configuration
	DBType       string          `json:"db_type" yaml:"DB_TYPE"`               // Database type (e.g., postgres)
	HostBasePath string          `json:"host_base_path" yaml:"HOST_BASE_PATH"` // Base host URL for Swagger/docs
	Metrics      MetricsSettings `json:"metrics" yaml:"METRICS"`               // Metrics server settings
	Port         int             `json:"port" yaml:"PORT"`                     // Application listening port
}

// MetricsSettings holds Prometheus metrics configuration.
type MetricsSettings struct {
	Path     string `envconfig:"METRICS_PATH"`     // Metrics endpoint path (e.g., /metrics)
	UserName string `envconfig:"METRICS_USERNAME"` // Optional basic auth username
	Password string `envconfig:"METRICS_PASSWORD"` // Optional basic auth password
	Port     int    `envconfig:"METRICS_PORT"`     // Metrics server port
}

// RedisConfig holds Redis connection details.
type RedisConfig struct {
	Host     string `json:"host" yaml:"HOST"`         // Redis host
	Port     string `json:"port" yaml:"PORT"`         // Redis port
	Password string `json:"password" yaml:"PASSWORD"` // Redis password
	TTL      string `json:"ttl" yaml:"TTL"`           // Default TTL for cached items
}

// LoadConfig loads application configuration from a YAML file and environment variables.
// Env variables take precedence over YAML values if both are provided.
func LoadConfig(filePath string) (*Config, error) {
	cfg := Config{}

	// Load from YAML file if provided
	if filePath != "" {
		if err := readFile(&cfg, filePath); err != nil {
			return nil, err
		}
	}

	// Override with environment variables
	if err := readEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// readFile reads configuration from a YAML file and decodes it into Config struct.
func readFile(cfg *Config, filePath string) error {
	if filePath == "" {
		// No file path provided, skip file reading
		return nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(cfg); err != nil {
		return err
	}

	return nil
}

// readEnv processes environment variables using envconfig and overrides existing config.
func readEnv(cfg *Config) error {
	// Empty prefix "" means no prefix is required for env variables
	return envconfig.Process("", cfg)
}
