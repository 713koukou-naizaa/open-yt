package youtube

import (
	"encoding/json"
	"fmt"
)

// Fetches latest videos from user's YouTube home feed
// Requires yt-dlp to be configured with cookies for logged-in YouTube account
func HomeFeed(paginationThreshold int, YTDLPCommand string, browser string) ([]YTVideo, error) {
	// Command to fetch latest home feed videos as JSON
	// "youtube:home" is a special playlist for yt-dlp
	YTDLPArgs := []string{
		"--flat-playlist",
		"--dump-json",
		"--playlist-end",
		fmt.Sprintf("%d", paginationThreshold)}
	if browser != "" {
		YTDLPArgs = append(YTDLPArgs, "--cookies-from-browser", browser)
	}
	YTDLPArgs = append(YTDLPArgs, BaseURL)

	var YTVideos []YTVideo
	processor := func(line []byte) error {
		var currentYTDLPVideo YTDLPVideo
		if err := json.Unmarshal(line, &currentYTDLPVideo); err != nil {
			return err // Indicates a non-JSON line
		}
		YTVideos = append(YTVideos, newYTVideoFromYTDLPVideo(currentYTDLPVideo))
		return nil
	}

	if err := ytdlpExecutor(YTDLPCommand, YTDLPArgs, processor); err != nil {
		return nil, fmt.Errorf("home feed fetch failed: %w", err)
	}

	return YTVideos, nil
}