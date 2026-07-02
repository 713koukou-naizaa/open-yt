package youtube

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
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

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start yt-dlp command for subscriptions: %w", err)
	}

	var YTVideos []YTVideo
	// Read stderr first to ensure we capture any error messages if the command fails early.
	// We can't read it after cmd.Wait() because the pipe will be closed.
	stderrBytes, _ := io.ReadAll(stderr)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var currentYTDLPVideo YTDLPVideo
		line := scanner.Bytes()

		currentYTDLPVideoUnmarshallError := json.Unmarshal(line, &currentYTDLPVideo)
		if currentYTDLPVideoUnmarshallError != nil {
			// It's possible yt-dlp outputs non-json lines (e.g., warnings about cookies)
			fmt.Fprintf(os.Stderr, "yt-dlp warning: %s\n", string(line))
			continue
		}

		YTVideos = append(YTVideos, newYTVideoFromYTDLPVideo(currentYTDLPVideo))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading yt-dlp subscription output: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		// Provide more informative error message including yt-dlp's output
		return nil, fmt.Errorf("yt-dlp command for subscriptions failed: %w\n%s", err, string(stderrBytes))
	}

	return YTVideos, nil
}
