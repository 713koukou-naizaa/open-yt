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
		case cmdSearch:
			return a.runSearch(args[1:])

		case cmdPlay:
			return a.runPlay(args[1:])

		case cmdHelp, "-h", "--help":
			return a.printHelp()

		default:
			return fmt.Errorf("unknown command: %s", args[0])
	}
}

func (a App) runInteractive() error {
	p := tea.NewProgram(NewMainMenuModel())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running interactive menu: %w", err)
	}

	m := finalModel.(SimpleMenuModel)

	// Create new scanner to read from stdin
	scanner := bufio.NewScanner(os.Stdin)

	switch m.selected {
		case menuSearch:
			fmt.Print("Enter search query: ")
			if scanner.Scan() {
				query := scanner.Text()
				return a.runSearch(strings.Fields(query))
			}

		case menuPlay:
			fmt.Print("Enter YouTube URL or video ID: ")
			if scanner.Scan() {
				url := scanner.Text()
				return a.runPlay([]string{url})
			}

		case menuSubscriptions:
			return a.runSubscriptions()

		case menuSubscriptionsFeed:
			return a.runSubscriptionsFeed()

		case menuQuit:
			return nil
		}
	return nil
}

func (a App) runSubscriptions() error {
	// Fetch channels
	fmt.Println("Fetching your subscriptions...")
	channels, err := youtube.GetSubscriptions(a.config.YTDLPCommand, a.config.CookiesFromBrowser)
	if err != nil {
		return err
	}
	if len(channels) == 0 {
		fmt.Println("Could not find any subscriptions.")
		fmt.Println("Please ensure you have configured yt-dlp with cookies for a logged-in YouTube account.")
		return nil
	}

	channelNames := make([]string, len(channels))
	channelMap := make(map[string]youtube.YTChannel)
	for i, currentChannel := range channels {
		channelNames[i] = currentChannel.Name
		channelMap[currentChannel.Name] = currentChannel
	}

	// Select channel
	listModel := NewFilterableListModel(channelNames, "Select a channel:")
	p := tea.NewProgram(listModel)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running channel selection: %w", err)
	}

	selectedChanModel := finalModel.(FilterableListModel)
	if selectedChanModel.back || selectedChanModel.selected == "" {
		return a.runInteractive() // Go back to main menu
	}

	selectedChannel := channelMap[selectedChanModel.selected]

	// Select channel content type
	contentTypeModel := NewSimpleMenuModel([]string{contentTypeVideos, contentTypeStreams}, "Select content type:")
	p = tea.NewProgram(contentTypeModel)
	finalModel, err = p.Run()
	if err != nil {
		return fmt.Errorf("error running content type selection: %w", err)
	}

	selectedContentTypeModel := finalModel.(SimpleMenuModel)
	if selectedContentTypeModel.selected == "" { 
		return a.runInteractive()
	}

	var contentType string
	if selectedContentTypeModel.selected == contentTypeVideos {
		contentType = contentTypeVideosInternal
	} else {
		contentType = contentTypeStreamsInternal
	}

	// Fetch channel content type content
	fmt.Printf("Fetching latest %s from %s...\n", contentType, selectedChannel.Name)
	videos, err := youtube.GetChannelUploads(selectedChannel.ID, contentType, a.config.PaginationThreshold, a.config.YTDLPCommand)
	if err != nil {
		return err
	}
	if len(videos) == 0 {
		fmt.Printf("No recent %s found for %s.\n", contentType, selectedChannel.Name)
		return nil
	}

	// Select channel content type content
	return a.runInteractiveVideoList(videos)
}

func (a App) runSearch(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: open-yt %s <query>", cmdSearch)
	}

	// Fetch videos
	query := strings.Join(args, " ")
	videos, err := a.searchVideos(query)
	if err != nil || videos == nil {
		return err
	}

	// Select video
	return a.runInteractiveVideoList(videos)
}

func (a App) runSubscriptionsFeed() error {
	// Fetch videos
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

	// Select video
	return a.runInteractiveVideoList(videos)
}

func (a App) runPlay(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: open-yt %s <youtube-url-or-video-id>", cmdPlay)
	}

	// Fetch and play video
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
		return a.runInteractive()
	}
	return nil
}

// Helper to abstract video searching logic
// Returns a slice of videos or a nil slice if no videos were found
func (a App) searchVideos(query string) ([]youtube.YTVideo, error) {
	// Fetch videos
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
	// Fetch and play video
	return player.Play(videoURL, a.config.Player)
}

func (a App) printHelp() error {
	fmt.Println(`open-yt

Usage:
  open-yt ` + cmdSearch + ` <query>
  open-yt ` + cmdPlay + ` <youtube-url-or-video-id>

Commands:
  ` + cmdSearch + `    Search for videos
  ` + cmdPlay + `      Play a video with the configured player
  ` + cmdHelp + `      Show this help message`)

	return nil
}
