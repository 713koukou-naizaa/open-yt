package youtube

import (
	"encoding/json"
	"fmt"
)

// Fetches the list of channels the user is subscribed to.
func GetSubscriptions(YTDLPCommand string, YTDLPCookiesBrowser string) ([]YTChannel, error) {
	YTDLPArgs := []string{"--flat-playlist", "--dump-json"}
	if YTDLPCookiesBrowser != "" {
		YTDLPArgs = append(YTDLPArgs, "--cookies-from-browser", YTDLPCookiesBrowser)
	}
	YTDLPArgs = append(YTDLPArgs, feedChannelsURL)

	var YTChannels []YTChannel
	processor := func(line []byte) error {
		var currentYTChannel struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			URL   string `json:"url"`
		}
		if err := json.Unmarshal(line, &currentYTChannel); err != nil {
			return err
		}
		YTChannels = append(YTChannels, YTChannel{
			ID:   currentYTChannel.ID,
			Name: currentYTChannel.Title,
			URL:  currentYTChannel.URL,
		})
		return nil
	}

	if err := ytdlpExecutor(YTDLPCommand, YTDLPArgs, processor); err != nil {
		return nil, fmt.Errorf("subscriptions fetch failed: %w", err)
	}

	return YTChannels, nil
}

// Fetches videos from a specific channel's tab
func GetChannelUploads(channelID, contentType string, paginationThreshold int, YTDLPCommand string) ([]YTVideo, error) {
	channelURLString := fmt.Sprintf(channelURLFormat, channelID, contentType)

	YTDLPArgs := []string{
		"--flat-playlist",
		"--dump-json",
		"--playlist-end",
		fmt.Sprintf("%d", paginationThreshold),
		channelURLString,
	}

	var YTvideos []YTVideo
	processor := func(line []byte) error {
		var currentYTDLPVideo YTDLPVideo
		if err := json.Unmarshal(line, &currentYTDLPVideo); err != nil {
			return err
		}
		YTvideos = append(YTvideos, newYTVideoFromYTDLPVideo(currentYTDLPVideo))
		return nil
	}

	if err := ytdlpExecutor(YTDLPCommand, YTDLPArgs, processor); err != nil {
		return nil, fmt.Errorf("channel uploads fetch failed: %w", err)
	}

	return YTvideos, nil
}
