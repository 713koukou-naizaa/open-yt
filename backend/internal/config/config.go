package config

type Config struct {
	PlayerCommand string
	YTDLPCommand string
	MaxResults int
}

func Default() Config {
	return Config{
		PlayerCommand: "mpv",
		YTDLPCommand: "yt-dlp",
		MaxResults: 10,
	}
}
