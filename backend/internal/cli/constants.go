package cli

// Main Menu Options
const (
	// Display values for the main menu
	menuHomeFeed          = "Home feed"
	menuSubscriptions     = "Subscriptions"
	menuSubscriptionsFeed = "Subscriptions feed"
	menuSearch            = "Search"
	menuPlay              = "Play"
	menuQuit              = "Quit"
)

// Subscription Content Types
const (
	// Display values
	contentTypeVideosDisplay  = "Videos"
	contentTypeStreamsDisplay = "Streams"

	// Internal values used with yt-dlp
	contentTypeVideosLink  = "videos"
	contentTypeStreamsLink = "streams"
)

// CLI commands
const (
	cmdSearch = "search"
	cmdPlay   = "play"
	cmdHelp   = "help"
)
