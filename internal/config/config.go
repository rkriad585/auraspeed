package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	cfg         *Config
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

	configDir = filepath.Join(home, "."+appName)
	dataDir = filepath.Join(configDir, "data")
	configFile = filepath.Join(configDir, "config.toml")
	historyFile = filepath.Join(dataDir, "history.json")

	for _, dir := range []string{configDir, dataDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
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
		cfg = &defaultConfig
		return nil
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
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

// EnsureSensitiveFilePermissions ensures all sensitive files have restricted permissions (0600).
func EnsureSensitiveFilePermissions() error {
	files := []struct {
		path string
		perm os.FileMode
	}{
		{GetConfigFile(), 0600},
		{GetHistoryFile(), 0600},
		{filepath.Join(GetDataDir(), "servers.json"), 0600},
	}

	for _, f := range files {
		if err := ensureFilePermissions(f.path, f.perm); err != nil {
			return fmt.Errorf("failed to set permissions on %s: %w", f.path, err)
		}
	}
	return nil
}
