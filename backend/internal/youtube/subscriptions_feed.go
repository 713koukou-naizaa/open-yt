package youtube

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
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

	cmd := exec.Command(YTDLPCommand, YTDLPArgs...)

	stdout, cmdStdoutPipeError := cmd.StdoutPipe()
	if cmdStdoutPipeError != nil {
		return nil, fmt.Errorf("failed to get stdout pipe for subscriptions: %w", cmdStdoutPipeError)
	}

	// Capture stderr to provide better error messages
	stderr, cmdStderrPipeError := cmd.StderrPipe()
	if cmdStderrPipeError != nil {
		return nil, fmt.Errorf("failed to get stderr pipe for subscriptions: %w", cmdStderrPipeError)
	}

	cmdStdStartError := cmd.Start()
	if cmdStdStartError != nil {
		return nil, fmt.Errorf("failed to start yt-dlp command for subscriptions: %w", cmdStdStartError)
	}

	var YTVideos []YTVideo
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var currentYTDLPVideo YTDLPVideo
		line := scanner.Bytes()

		currentYTDLPVideoUnmarshallError := json.Unmarshal(line, &currentYTDLPVideo)
		if currentYTDLPVideoUnmarshallError != nil {
			// It's possible yt-dlp outputs non-json lines (e.g., warnings about cookies)
			// For now, continue
			continue
		}

		YTVideos = append(YTVideos, newYTVideoFromYTDLPVideo(currentYTDLPVideo))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading yt-dlp subscription output: %w", err)
	}

	cmdWaitError := cmd.Wait()
	if cmdWaitError != nil {
		stderrBytes, _ := io.ReadAll(stderr)
		// Provide more informative error message including yt-dlp's output
		return nil, fmt.Errorf("yt-dlp command for subscriptions failed: %w\n%s", cmdWaitError, string(stderrBytes))
	}

	return YTVideos, nil
}
