package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// Config holds all configuration for the application.
// Tags are used by viper to map configuration keys to struct fields.
type Config struct {
	YTDLPCommand        string `mapstructure:"yt-dlp_command"`
	Player              Player `mapstructure:"player"`
	CookiesFromBrowser  string `mapstructure:"cookies_from_browser"`
	PaginationThreshold int    `mapstructure:"pagination_threshold"`
}

// Player holds the configuration for the video player.
type Player struct {
	Command         string `mapstructure:"command"`
	YTDLFormat      string `mapstructure:"ytdl_format"`
	Volume          *int   `mapstructure:"volume"` // Pointer to distinguish between 0 and not set
	Fullscreen      string `mapstructure:"fullscreen"`
	WindowMaximized string `mapstructure:"window_maximized"`
	KeepOpen        string `mapstructure:"keep_open"`
	ForceWindow     string `mapstructure:"force_window"`
}

// Load reads configuration from a file and sets defaults.
func Load() (Config, error) {
	var cfg Config

	// Set default values for top-level config
	viper.SetDefault("yt-dlp_command", "yt-dlp")
	viper.SetDefault("player.command", "mpv") // Default player command
	viper.SetDefault("player.ytdl_format", "")
	viper.SetDefault("player.volume", 100)
	viper.SetDefault("player.fullscreen", "yes")
	viper.SetDefault("player.window_maximized", "yes")
	viper.SetDefault("player.keep_open", "yes")
	viper.SetDefault("player.force_window", "yes")
	viper.SetDefault("cookies_from_browser", "firefox")
	viper.SetDefault("pagination_threshold", 30)

	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, fmt.Errorf("could not get user home directory: %w", err)
	}

	// Set config path
	configPath := filepath.Join(home, ".config", "open-yt")
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			return cfg, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal the config into the struct
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return cfg, nil
}