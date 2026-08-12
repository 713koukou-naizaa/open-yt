package youtube

import (
	"encoding/json"
	"fmt"
)

// Fetches latest videos from user's subscriptions feed
// Requires yt-dlp to be configured with cookies for logged-in YouTube account
func SubscriptionsFeed(paginationThreshold int, YTDLPCommand string, browser string) ([]YTVideo, error) {
	// Command to fetch latest subscription videos as JSON
	YTDLPArgs := []string{
		"--flat-playlist",
		"--dump-json",
		"--playlist-end",
		fmt.Sprintf("%d", paginationThreshold)}
	if browser != "" {
		YTDLPArgs = append(YTDLPArgs, "--cookies-from-browser", browser)
	}
	YTDLPArgs = append(YTDLPArgs, feedSubscriptionsURL)

	var YTVideos []YTVideo
	processor := func(line []byte) error {
		var currentYTDLPVideo YTDLPVideo
		if err := json.Unmarshal(line, &currentYTDLPVideo); err != nil {
			return err
		}
		YTVideos = append(YTVideos, newYTVideoFromYTDLPVideo(currentYTDLPVideo))
		return nil
	}

	if err := ytdlpExecutor(YTDLPCommand, YTDLPArgs, processor); err != nil {
		return nil, fmt.Errorf("subscriptions feed fetch failed: %w", err)
	}

	return YTVideos, nil
}
