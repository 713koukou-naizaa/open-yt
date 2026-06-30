package config

// Holds configuration for media player
type Player struct {
	Command         string
	YTDLFormat      string // e.g., 'bestvideo[height<=720]+bestaudio/best[height<=720]'
	Volume          *int   // Pointer to distinguish between 0 and not set
	Fullscreen      string // "yes", "no", or "" for default
	WindowMaximized string // "yes", "no", or "" for default
	KeepOpen        string // "yes", "no", or "" for default
	ForceWindow     string // "yes", "no", or "" for default
}

type Config struct {
	Player       Player
	YTDLPCommand string
	MaxResults   int
	CookiesFromBrowser string
}

func Default() Config {
	defaultVolume := 35
	return Config{
		Player: Player{
			Command:     "mpv",
			KeepOpen:    "yes",
			ForceWindow: "yes",
			Volume:      &defaultVolume,
			Fullscreen:  "yes",
			WindowMaximized: "yes",
			YTDLFormat:  "bestvideo[height<=720]+bestaudio/best[height<=720]",
		},
		YTDLPCommand: "yt-dlp",
		MaxResults:   10,
		CookiesFromBrowser: "brave", // e.g., "firefox", "chrome", "brave"
	}
}
