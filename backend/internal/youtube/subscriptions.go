package youtube

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

// Fetches the list of channels the user is subscribed to.
func GetSubscriptions(YTDLPCommand string, YTDLPCookiesBrowser string) ([]YTChannel, error) {
	YTDLPArgs := []string{"--flat-playlist", "--dump-json"}
	if YTDLPCookiesBrowser != "" {
		YTDLPArgs = append(YTDLPArgs, "--cookies-from-browser", YTDLPCookiesBrowser)
	}
	YTDLPArgs = append(YTDLPArgs, feedChannelsURL)

	cmd := exec.Command(YTDLPCommand, YTDLPArgs...)

	stdout, cmdStdoutPipeError := cmd.StdoutPipe()
	if cmdStdoutPipeError != nil {
		return nil, fmt.Errorf("failed to get stdout pipe for subscriptions: %w", cmdStdoutPipeError)
	}

	stderr, cmdStderrPipeError := cmd.StderrPipe()
	if cmdStderrPipeError != nil {
		return nil, fmt.Errorf("failed to get stderr pipe for subscriptions: %w", cmdStderrPipeError)
	}

	cmdStartError := cmd.Start()
	if cmdStartError != nil {
		return nil, fmt.Errorf("failed to start yt-dlp command for subscriptions: %w", cmdStartError)
	}

	var YTChannels []YTChannel
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var currentYTChannel struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			URL   string `json:"url"`
		}
		line := scanner.Bytes()

		currentYTChannelUnmarshallError := json.Unmarshal(line, &currentYTChannel)
		if currentYTChannelUnmarshallError != nil {
			continue // Ignore non-JSON lines
		}
		YTChannels = append(YTChannels, YTChannel{
			ID:   currentYTChannel.ID,
			Name: currentYTChannel.Title,
			URL:  currentYTChannel.URL,
		})
	}

	scannerError := scanner.Err()
	if scannerError != nil {
		return nil, fmt.Errorf("error reading yt-dlp subscription output: %w", scannerError)
	}

	waitError := cmd.Wait()
	if waitError != nil {
		stderrBytes, _ := io.ReadAll(stderr)
		return nil, fmt.Errorf("yt-dlp command for subscriptions failed: %w\n%s", waitError, string(stderrBytes))
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

	cmd := exec.Command(YTDLPCommand, YTDLPArgs...)
	stdout, cmdStdoutPipeError := cmd.StdoutPipe()
	if cmdStdoutPipeError != nil {
		return nil, fmt.Errorf("failed to get stdout pipe for channel uploads: %w", cmdStdoutPipeError)
	}

	cmdStartError := cmd.Start()
	if cmdStartError != nil {
		return nil, fmt.Errorf("failed to start yt-dlp command for channel uploads: %w", cmdStartError)
	}

	var YTvideos []YTVideo
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var currentYTDLPVideo YTDLPVideo
		YTDLPVideoUnmarshallError := json.Unmarshal(scanner.Bytes(), &currentYTDLPVideo)
		if YTDLPVideoUnmarshallError != nil {
			continue
		}
		YTvideos = append(YTvideos, newYTVideoFromYTDLPVideo(currentYTDLPVideo))
	}

	cmd.Wait() // Wait for command to finish, ignore error as it might fail if playlist is empty

	return YTvideos, nil
}
