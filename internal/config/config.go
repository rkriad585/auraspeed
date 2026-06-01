package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var (
	cfg         *Config
	cfgMu       sync.RWMutex
	configDir   string
	configFile  string
	historyFile string
	dataDir     string
)

// Config holds all configuration for AuraSpeed.
type Config struct {
	Global    GlobalConfig      `mapstructure:"global"`
	Speedtest SpeedtestConfig   `mapstructure:"speedtest"`
	UI        UIConfig          `mapstructure:"ui"`
	Aliases   map[string]string `mapstructure:"aliases"`
}

// GlobalConfig holds global application settings.
type GlobalConfig struct {
	LogLevel    string `mapstructure:"loglevel"`
	NoColor     bool   `mapstructure:"nocolor"`
	AutoUpdate  bool   `mapstructure:"autoupdate"`
	ConfirmExit bool   `mapstructure:"confirmexit"`
}

// SpeedtestConfig holds speed test related settings.
type SpeedtestConfig struct {
	DefaultServerID   int `mapstructure:"defaultserverid"`
	Timeout           int `mapstructure:"timeout"`
	ParallelDownloads int `mapstructure:"paralleldownloads"`
	ParallelUploads   int `mapstructure:"paralleluploads"`
}

// UIConfig holds terminal UI related settings.
type UIConfig struct {
	Theme        string `mapstructure:"theme"`
	GraphHeight  int    `mapstructure:"graphheight"`
	HistoryLimit int    `mapstructure:"historylimit"`
	AutoRefresh  bool   `mapstructure:"autorefresh"`
	RefreshRate  int    `mapstructure:"refreshrate"`
	SaveHistory  bool   `mapstructure:"savehistory"`
}

// getHomeDir returns the user's home directory across platforms.
func getHomeDir() (string, error) {
	home := os.Getenv("USERPROFILE")
	if home != "" {
		return home, nil
	}
	home = os.Getenv("HOME")
	if home != "" {
		return home, nil
	}
	return os.UserHomeDir()
}

// Init initializes the configuration for the application.
// It sets up config directory, data directory, and loads config file.
// appName is used to determine the config directory name (e.g., ".auraspeed").
func Init(appName string) error {
	home, err := getHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir = filepath.Join(home, ".config", "neostore", appName)
	dataDir = filepath.Join(configDir, "data")
	configFile = filepath.Join(configDir, "config.toml")
	historyFile = filepath.Join(dataDir, "history.json")

	for _, dir := range []string{configDir, dataDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Set up environment variable support with AS_ prefix
	viper.SetEnvPrefix("AS")
	viper.AutomaticEnv()
	viper.BindEnv("loglevel")
	viper.BindEnv("nocolor")
	viper.BindEnv("autoupdate")
	viper.BindEnv("speedtest_timeout")
	viper.BindEnv("speedtest_defaultserverid")
	viper.BindEnv("ui_savehistory")
	viper.BindEnv("ui_historylimit")

	// Allow config file path to be overridden via env var
	if customPath := os.Getenv("AS_CONFIG_PATH"); customPath != "" {
		configFile = customPath
	}

	viper.SetConfigType("toml")
	viper.SetConfigFile(configFile)
	viper.SetDefault("global.loglevel", "info")
	viper.SetDefault("global.nocolor", false)
	viper.SetDefault("global.autoupdate", true)
	viper.SetDefault("global.confirmexit", false)
	viper.SetDefault("speedtest.timeout", 30)
	viper.SetDefault("speedtest.defaultserverid", 0)
	viper.SetDefault("speedtest.paralleldownloads", 4)
	viper.SetDefault("speedtest.paralleluploads", 2)
	viper.SetDefault("ui.theme", "default")
	viper.SetDefault("ui.graphheight", 8)
	viper.SetDefault("ui.historylimit", 100)
	viper.SetDefault("ui.autorefresh", false)
	viper.SetDefault("ui.refreshrate", 5)
	viper.SetDefault("ui.savehistory", true)

	if _, err := os.Stat(configFile); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}
	} else if os.IsNotExist(err) {
		defaultConfig := GetDefaultConfig()
		if err := viper.SafeWriteConfig(); err != nil {
			if err := viper.WriteConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Could not write default config: %v\n", err)
			}
		}
		cfgMu.Lock()
		cfg = &defaultConfig
		cfgMu.Unlock()
		return nil
	}

	var parsed Config
	if err := viper.Unmarshal(&parsed); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := parsed.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	cfgMu.Lock()
	cfg = &parsed
	cfgMu.Unlock()

	return nil
}

// Validate checks if the configuration values are valid.
func (c *Config) Validate() error {
	// Validate log level
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.Global.LogLevel] {
		return fmt.Errorf("invalid log level '%s', must be one of: debug, info, warn, error", c.Global.LogLevel)
	}

	// Validate timeout
	if c.Speedtest.Timeout < 10 || c.Speedtest.Timeout > 300 {
		return fmt.Errorf("invalid timeout '%d', must be between 10 and 300 seconds", c.Speedtest.Timeout)
	}

	// Validate parallel connections
	if c.Speedtest.ParallelDownloads < 1 || c.Speedtest.ParallelDownloads > 16 {
		return fmt.Errorf("invalid parallel downloads '%d', must be between 1 and 16", c.Speedtest.ParallelDownloads)
	}
	if c.Speedtest.ParallelUploads < 1 || c.Speedtest.ParallelUploads > 16 {
		return fmt.Errorf("invalid parallel uploads '%d', must be between 1 and 16", c.Speedtest.ParallelUploads)
	}

	// Validate UI settings
	if c.UI.GraphHeight < 3 || c.UI.GraphHeight > 20 {
		return fmt.Errorf("invalid graph height '%d', must be between 3 and 20", c.UI.GraphHeight)
	}
	if c.UI.HistoryLimit < 10 || c.UI.HistoryLimit > 1000 {
		return fmt.Errorf("invalid history limit '%d', must be between 10 and 1000", c.UI.HistoryLimit)
	}
	if c.UI.RefreshRate < 1 || c.UI.RefreshRate > 60 {
		return fmt.Errorf("invalid refresh rate '%d', must be between 1 and 60 seconds", c.UI.RefreshRate)
	}

	return nil
}

// GetDefaultConfig returns the default configuration values.
func GetDefaultConfig() Config {
	return Config{
		Global: GlobalConfig{
			LogLevel:    "info",
			NoColor:     false,
			AutoUpdate:  true,
			ConfirmExit: false,
		},
		Speedtest: SpeedtestConfig{
			DefaultServerID:   0,
			Timeout:           30,
			ParallelDownloads: 4,
			ParallelUploads:   2,
		},
		UI: UIConfig{
			Theme:        "default",
			GraphHeight:  8,
			HistoryLimit: 100,
			AutoRefresh:  false,
			RefreshRate:  5,
			SaveHistory:  true,
		},
		Aliases: map[string]string{
			"st":   "speedtest",
			"si":   "info",
			"net":  "network",
			"hist": "history",
		},
	}
}

// Get returns the current configuration.
// If not initialized, it returns a default config.
func Get() *Config {
	cfgMu.RLock()
	if cfg != nil {
		defer cfgMu.RUnlock()
		return cfg
	}
	cfgMu.RUnlock()

	cfgMu.Lock()
	defer cfgMu.Unlock()
	if cfg == nil {
		defaultCfg := GetDefaultConfig()
		cfg = &defaultCfg
	}
	return cfg
}

// GetConfigDir returns the configuration directory path.
func GetConfigDir() string {
	return configDir
}

// GetDataDir returns the data directory path.
func GetDataDir() string {
	return dataDir
}

// GetHistoryFile returns the history file path.
func GetHistoryFile() string {
	return historyFile
}

// GetConfigFile returns the config file path.
func GetConfigFile() string {
	return configFile
}

// ensureFilePermissions checks if a file has the expected permissions and fixes them if not.
// Returns nil if the file doesn't exist yet.
// ensureFilePermissions ensures the file at path has the specified permissions.
func ensureFilePermissions(path string, perm os.FileMode) error {
	info, err := os.Stat(path)
	if err != nil {
		return nil // file doesn't exist yet
	}
	if info.Mode().Perm() != perm {
		return os.Chmod(path, perm)
	}
	return nil
}

// EnsureSensitiveFilePermissions ensures config files have restricted permissions (0600).
// History file uses 0644 so users can read their test history.
func EnsureSensitiveFilePermissions() error {
	files := []struct {
		path string
		perm os.FileMode
	}{
		{GetConfigFile(), 0600},
		{GetHistoryFile(), 0644},
		{filepath.Join(GetDataDir(), "servers.json"), 0600},
	}

	for _, f := range files {
		if err := ensureFilePermissions(f.path, f.perm); err != nil {
			return fmt.Errorf("failed to set permissions on %s: %w", f.path, err)
		}
	}
	return nil
}
