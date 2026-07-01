package cli

// Main Menu Options
const (
	// Display values for the main menu
	menuSearch            = "Search"
	menuPlay              = "Play"
	menuSubscriptions     = "Subscriptions"
	menuSubscriptionsFeed = "Subscriptions feed"
	menuQuit              = "Quit"
)

// Subscription Content Types
const (
	// Display values
	contentTypeVideos  = "Videos"
	contentTypeStreams = "Streams"

	// Internal values used with yt-dlp
	contentTypeVideosInternal  = "videos"
	contentTypeStreamsInternal = "streams"
)

// CLI commands
const (
	cmdSearch = "search"
	cmdPlay   = "play"
	cmdHelp   = "help"
)