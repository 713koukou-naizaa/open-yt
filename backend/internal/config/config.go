package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Holds all configuration for the application
// Tags are used by viper to map configuration keys to struct fields
type AppConfiguration struct {
	YTDLPCommand        string              `mapstructure:"yt-dlp_command"`
	PlayerConfiguration PlayerConfiguration `mapstructure:"player"`
	CookiesFromBrowser  string              `mapstructure:"cookies_from_browser"`
	PaginationThreshold int                 `mapstructure:"pagination_threshold"`
}

// Player holds the configuration for the video player
type PlayerConfiguration struct {
	Command         string `mapstructure:"command"`
	YTDLFormat      string `mapstructure:"ytdl_format"`
	Volume          *int   `mapstructure:"volume"` // Pointer to distinguish between 0 and not set
	Fullscreen      string `mapstructure:"fullscreen"`
	WindowMaximized string `mapstructure:"window_maximized"`
	KeepOpen        string `mapstructure:"keep_open"`
	ForceWindow     string `mapstructure:"force_window"`
}

// Load reads configuration from a file and sets defaults
func Load() (AppConfiguration, error) {
	var appConfiguration AppConfiguration

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
	homeDirectory, getHomeDirectoryError := os.UserHomeDir()
	if getHomeDirectoryError != nil {
		return appConfiguration, fmt.Errorf("could not get user home directory: %w", getHomeDirectoryError)
	}

	// Set config path
	configurationFilePath := filepath.Join(homeDirectory, ".config", "open-yt")
	viper.AddConfigPath(configurationFilePath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Attempt to read the config file
	viperConfigurationFileReadError := viper.ReadInConfig()
	if viperConfigurationFileReadError != nil {
		_, viperConfigurationFileReadErrorIsFileNotFound := viperConfigurationFileReadError.(viper.ConfigFileNotFoundError)
		if !viperConfigurationFileReadErrorIsFileNotFound {
			// Configuration file was found but another error was produced
			return appConfiguration, fmt.Errorf("error reading config file: %w", viperConfigurationFileReadError)
		}
	}

	// Unmarshal the config into the struct
	viperConfigurationUnmarshallError := viper.Unmarshal(&appConfiguration)
	if viperConfigurationUnmarshallError != nil {
		return appConfiguration, fmt.Errorf("unable to decode into struct: %w", viperConfigurationUnmarshallError)
	}

	return appConfiguration, nil
}
