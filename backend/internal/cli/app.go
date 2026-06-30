package cli

import (
	"bufio"
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"strings"

	"open-yt/internal/config"
	"open-yt/internal/player"
	"open-yt/internal/youtube"
	"os"
)

type App struct {
	config config.Config
}

func NewApp(cfg config.Config) App {
	return App{config: cfg}
}

func (a App) Run(args []string) error {
	if len(args) == 0 {
		return a.runInteractive()
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

func (a App) runInteractive() error {
	p := tea.NewProgram(newMenuModel())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running interactive menu: %w", err)
	}

	m := finalModel.(menuModel)

	// Create a new scanner to read from stdin
	scanner := bufio.NewScanner(os.Stdin)

	switch m.selected {
	case "Search":
		fmt.Print("Enter search query: ")
		if scanner.Scan() {
			query := scanner.Text()
			return a.runSearch(strings.Fields(query))
		}
	case "Play":
		fmt.Print("Enter YouTube URL or video ID: ")
		if scanner.Scan() {
			url := scanner.Text()
			return a.runPlay([]string{url})
		}
	case "Subscriptions feed":
		return a.runSubscriptionsFeed()
	case "Quit":
		return nil
	}
	return nil
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

	return a.runInteractiveVideoList(videos)
}

func (a App) runSubscriptionsFeed() error {
	fmt.Println("Fetching videos from your subscriptions feed...")
	videos, err := youtube.SubscriptionsFeed(a.config.PaginationThreshold, a.config.YTDLPCommand, a.config.CookiesFromBrowser)
	if err != nil {
		return err
	}
	if len(videos) == 0 {
		fmt.Println("No videos found in your subscriptions feed.")
		fmt.Println("Please ensure you have configured yt-dlp with cookies for YouTube.")
		return nil
	}

	return a.runInteractiveVideoList(videos)
}

func (a App) runPlay(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: open-yt play <youtube-url-or-video-id>")
	}

	video := args[0]

	return a.playVideo(video)
}

// Displays list of videos and handles user selection
func (a App) runInteractiveVideoList(videos []youtube.YTVideo) error {
	p := tea.NewProgram(newYTSearchModel(videos))
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running interactive video list: %w", err)
	}

	m := finalModel.(YTSearchModel)

	// Delegate to result handler
	return a.handleSearchResult(m)
}

// Processes result from YTSearchModel bubble
func (a App) handleSearchResult(m YTSearchModel) error {
	// If user chose to go back, re-run main interactive menu
	if m.back {
		return a.runInteractive()
	}

	if m.selectedVideo != nil {
		err := a.playVideo(m.selectedVideo.URL)
		if err != nil {
			return err
		}
		return a.runInteractive() // Return to main menu after video finishes
	}
	return nil
}

// Helper to abstract video searching logic
// Returns a slice of videos or a nil slice if no videos were found
func (a App) searchVideos(query string) ([]youtube.YTVideo, error) {
	videos, err := youtube.Search(query, a.config.PaginationThreshold, a.config.YTDLPCommand)
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
