package cli

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
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

	videos, err := a.searchVideos(query)
	if err != nil || videos == nil { // videos == nil means no results were found and message was printed
		return err
	}
	
	for i, video := range videos {
		durationMinutes := int(video.Duration / 60)
		durationSeconds := int(video.Duration) % 60
		fmt.Printf("%d. %s\n", i+1, video.Title)
		fmt.Printf("   Channel: %s | Duration: %d:%02d | URL: %s\n", video.Channel, durationMinutes, durationSeconds, video.URL)
	}

	p := tea.NewProgram(newYTSearchModel(videos))
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running interactive search: %w", err)
	}

	// State of model when program exits
	m := finalModel.(YTSearchModel)

	if m.selectedVideo != nil {
		return a.playVideo(m.selectedVideo.URL)
	}

	return nil
}

func (a App) runPlay(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: open-yt play <youtube-url-or-video-id>")
	}

	video := args[0]

	return a.playVideo(video)
}

// Helper to abstract video searching logic
// Returns a slice of videos or a nil slice if no videos were found
func (a App) searchVideos(query string) ([]youtube.YTVideo, error) {
	videos, err := youtube.Search(query, a.config.MaxResults, a.config.YTDLPCommand)
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		fmt.Println("No videos found for your query.")
		return nil, nil
	}
	return videos, nil
}

// Helper to abstract the video playing logic
func (a App) playVideo(videoURL string) error {
	return player.Play(videoURL, a.config.Player)
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
