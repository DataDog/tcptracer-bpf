package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	DefaultLogFilePath    = "/var/log/datadog/network-tracer.log"
	DefaultUnixSocketPath = "/opt/datadog-agent/run/nettracer.sock"
)

// Config is the global config for the network-tracer. This information is sourced from config files and
// the environment variables.
type Config struct {
	Enabled        bool
	LogFile        string
	LogLevel       string
	LogToConsole   bool
	UnixSocketPath string
}

// NewDefaultConfig returns an Config with defaults initialized
func NewDefaultConfig() *Config {
	return &Config{
		Enabled:        false,
		LogFile:        DefaultLogFilePath,
		LogLevel:       "info",
		LogToConsole:   false,
		UnixSocketPath: DefaultUnixSocketPath,
	}
}

// NewConfig returns an Config using a configuration file. It can be nil
// if there is no file available. In this case we'll configure only via environment.
func NewConfig(yamlCfg *YamlConfig) (*Config, error) {
	var err error
	cfg := NewDefaultConfig()

	// For Agents >= 6 we will have a YAML config file to use.
	if yamlCfg != nil {
		if cfg, err = mergeYamlConfig(cfg, yamlCfg); err != nil {
			return nil, err
		}
	}

	// Use environment to override config values
	cfg = mergeEnvironmentVariables(cfg)

	// Python-style log level has WARNING vs WARN
	if strings.ToLower(cfg.LogLevel) == "warning" {
		cfg.LogLevel = "warn"
	}

	// (Re)configure the logging from our configuration
	if err := NewLoggerLevel(cfg.LogLevel, cfg.LogFile, cfg.LogToConsole); err != nil {
		return nil, err
	}

	return cfg, nil
}

// mergeEnvironmentVariables applies overrides from environment variables to the process agent configuration
func mergeEnvironmentVariables(cfg *Config) *Config {
	if enabled, err := isAffirmative(os.Getenv("DD_NETWORK_TRACING_ENABLED")); err == nil {
		cfg.Enabled = enabled
	}

	// Network tracer unix socket location
	if v := os.Getenv("DD_NETTRACER_SOCKET"); v != "" {
		cfg.UnixSocketPath = v
	}

	// Support LOG_LEVEL and DD_LOG_LEVEL, but prefer DD_LOG_LEVEL
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("DD_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}

	// Support DD_LOGS_STDOUT and LOG_TO_CONSOLE, but prefer LOG_TO_CONSOLE
	if enabled, err := isAffirmative(os.Getenv("DD_LOGS_STDOUT")); err == nil {
		cfg.LogToConsole = enabled
	}
	if enabled, err := isAffirmative(os.Getenv("LOG_TO_CONSOLE")); err == nil {
		cfg.LogToConsole = enabled
	}

	return cfg
}

func isAffirmative(value string) (bool, error) {
	if value == "" {
		return false, fmt.Errorf("value is empty")
	}
	v := strings.ToLower(value)
	return v == "true" || v == "yes" || v == "1", nil
}
