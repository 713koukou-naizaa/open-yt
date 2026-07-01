package youtube

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

// GetSubscriptions fetches the list of channels the user is subscribed to.
func GetSubscriptions(YTDLPCommand string, browser string) ([]YTChannel, error) {
	args := []string{"--flat-playlist", "--dump-json"}
	if browser != "" {
		args = append(args, "--cookies-from-browser", browser)
	}
	args = append(args, feedChannelsURL)

	cmd := exec.Command(YTDLPCommand, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe for subscriptions: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stderr pipe for subscriptions: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start yt-dlp command for subscriptions: %w", err)
	}

	var channels []YTChannel
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var currentChannel struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			URL   string `json:"url"`
		}
		line := scanner.Bytes()

		if err := json.Unmarshal(line, &currentChannel); err != nil {
			continue // Ignore non-JSON lines
		}
		channels = append(channels, YTChannel{
			ID:   currentChannel.ID,
			Name: currentChannel.Title,
			URL:  currentChannel.URL,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading yt-dlp subscription output: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		stderrBytes, _ := io.ReadAll(stderr)
		return nil, fmt.Errorf("yt-dlp command for subscriptions failed: %w\n%s", err, string(stderrBytes))
	}

	return channels, nil
}

// GetChannelUploads fetches videos from a specific channel's "videos" or "live" tab.
func GetChannelUploads(channelID, contentType string, paginationThreshold int, YTDLPCommand string) ([]YTVideo, error) {
	// contentType should be "videos" or "live"
	channelURL := fmt.Sprintf(channelURLFormat, channelID, contentType)

	args := []string{
		"--flat-playlist",
		"--dump-json",
		"--playlist-end",
		fmt.Sprintf("%d", paginationThreshold),
		channelURL,
	}

	cmd := exec.Command(YTDLPCommand, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe for channel uploads: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start yt-dlp command for channel uploads: %w", err)
	}

	var videos []YTVideo
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var currentYTDLPVideo YTDLPVideo
		if err := json.Unmarshal(scanner.Bytes(), &currentYTDLPVideo); err != nil {
			continue
		}
		videos = append(videos, newYTVideoFromYTDLPVideo(currentYTDLPVideo))
	}

	cmd.Wait() // Wait for the command to finish, ignore error as it might fail if playlist is empty

	return videos, nil
}