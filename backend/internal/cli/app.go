package cli

import (
	"bufio"
	"fmt"
	"strings"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"open-yt/internal/config"
	"open-yt/internal/player"
	"open-yt/internal/youtube"
)

type Application struct {
	configuration config.AppConfiguration
}

func NewApplication(configuration config.AppConfiguration) Application {
	return Application{configuration: configuration}
}

func (application Application) Run(applicationArgs []string) error {
	if len(applicationArgs) == 0 {
		return application.runInteractive()
	}

	switch applicationArgs[0] {
	case cmdHelp, "-h", "--help":
		return application.printHelp()

	default:
		return fmt.Errorf("unknown command: %s", applicationArgs[0])
	}
}

func (a Application) runInteractive() error {
	p := tea.NewProgram(NewMainMenuModel())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running interactive menu: %w", err)
	}

	m := finalModel.(SimpleMenuModel)

	// Create new scanner to read from stdin
	scanner := bufio.NewScanner(os.Stdin)

	switch m.selectedChoice {
	case menuHomeFeed:
		return a.runHomeFeed()

	case menuSubscriptions:
		return a.runSubscriptions()

	case menuSubscriptionsFeed:
		return a.runSubscriptionsFeed()

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

	case menuQuit:
		return nil
	}
	return nil
}

func (a Application) runHomeFeed() error {
	// Fetch videos
	fmt.Println("Fetching videos from your home feed...")
	videos, err := youtube.HomeFeed(a.configuration.PaginationThreshold, a.configuration.YTDLPCommand, a.configuration.CookiesFromBrowser)
	if err != nil {
		return err
	}
	if len(videos) == 0 {
		fmt.Println("No videos found in your home feed.")
		fmt.Println("Please ensure you have configured yt-dlp with cookies for YouTube.")
		return nil
	}

	// Select video
	return a.runInteractiveVideoList(videos)
}

func (a Application) runSubscriptions() error {
	// Fetch channels
	fmt.Println("Fetching your subscriptions...")
	channels, err := youtube.GetSubscriptions(a.configuration.YTDLPCommand, a.configuration.CookiesFromBrowser)
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
	if selectedChanModel.shouldGoBack || selectedChanModel.selectedChoice == "" {
		return a.runInteractive() // Go back to main menu
	}

	selectedChannel := channelMap[selectedChanModel.selectedChoice]

	// Select channel content type
	contentTypeModel := NewSimpleMenuModel([]string{contentTypeVideosDisplay, contentTypeStreamsDisplay}, "Select content type:")
	p = tea.NewProgram(contentTypeModel)
	finalModel, err = p.Run()
	if err != nil {
		return fmt.Errorf("error running content type selection: %w", err)
	}

	selectedContentTypeModel := finalModel.(SimpleMenuModel)
	if selectedContentTypeModel.selectedChoice == "" {
		return a.runInteractive()
	}

	var contentType string
	if selectedContentTypeModel.selectedChoice == contentTypeVideosDisplay {
		contentType = contentTypeVideosLink
	} else {
		contentType = contentTypeStreamsLink
	}

	// Fetch channel content type content
	fmt.Printf("Fetching latest %s from %s...\n", contentType, selectedChannel.Name)
	videos, err := youtube.GetChannelUploads(selectedChannel.ID, contentType, a.configuration.PaginationThreshold, a.configuration.YTDLPCommand)
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

func (a Application) runSubscriptionsFeed() error {
	// Fetch videos
	fmt.Println("Fetching videos from your subscriptions feed...")
	videos, err := youtube.SubscriptionsFeed(a.configuration.PaginationThreshold, a.configuration.YTDLPCommand, a.configuration.CookiesFromBrowser)
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

func (a Application) runSearch(searchArgs []string) error {
	if len(searchArgs) == 0 {
		return fmt.Errorf("usage: open-yt %s <query>", cmdSearch)
	}

	// Fetch videos
	query := strings.Join(searchArgs, " ")
	videos, err := a.searchVideos(query)
	if err != nil || videos == nil {
		return err
	}

	// Select video
	return a.runInteractiveVideoList(videos)
}

func (a Application) runPlay(playArgs []string) error {
	if len(playArgs) == 0 {
		return fmt.Errorf("usage: open-yt %s <youtube-url-or-video-id>", cmdPlay)
	}

	// Fetch and play video
	video := playArgs[0]
	return a.playVideo(video)
}

// Displays list of videos and handles user selection
func (a Application) runInteractiveVideoList(videos []youtube.YTVideo) error {
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
func (a Application) handleSearchResult(m YTSearchModel) error {
	// If user chose to go back, re-run main interactive menu
	if m.shouldGoBack {
		return a.runInteractive()
	}

	if m.selectedYTVideo != nil {
		err := a.playVideo(m.selectedYTVideo.URL)
		if err != nil {
			return err
		}
		return a.runInteractive()
	}
	return nil
}

// Helper to abstract video searching logic
// Returns a slice of videos or a nil slice if no videos were found
func (a Application) searchVideos(query string) ([]youtube.YTVideo, error) {
	// Fetch videos
	videos, err := youtube.Search(query, a.configuration.PaginationThreshold, a.configuration.YTDLPCommand)
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
func (a Application) playVideo(videoURL string) error {
	// Fetch and play video
	return player.Play(videoURL, a.configuration.PlayerConfiguration)
}

func (a Application) printHelp() error {
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
