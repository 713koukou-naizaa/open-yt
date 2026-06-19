package cli

import (
	"fmt"
	"strings"

	"open-yt/internal/config"
	"open-yt/internal/player"
	"open-yt/internal/youtube"
)

type App struct {
	config config.Config
}

func NewApp(cfg config.Config) App {
	return App{config: cfg}
}

func (a App) Run(args []string) error {
	if len(args) == 0 {
		return a.printHelp()
	}

	switch args[0] {
	case "search":
		return a.runSearch(args[1:])
	case "play":
		return a.runPlay(args[1:])
	case "help", "-h", "--help":
		return a.printHelp()
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func (a App) runSearch(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: open-yt search <query>")
	}

	query := strings.Join(args, " ")

	videos, err := youtube.Search(query, a.config.MaxResults)
	if err != nil {
		return err
	}

	for i, video := range videos {
		fmt.Printf("%d. %s\n", i+1, video.Title)
		fmt.Printf("   %s\n", video.URL)
	}

	return nil
}

func (a App) runPlay(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: open-yt play <youtube-url-or-video-id>")
	}

	video := args[0]

	return player.Play(video, a.config.PlayerCommand)
}

func (a App) printHelp() error {
	fmt.Println(`open-yt

Usage:
  open-yt search <query>
  open-yt play <youtube-url-or-video-id>

Commands:
  search    Search for videos
  play      Play a video with the configured player
  help      Show this help message`)

	return nil
}
