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
	args := []string{"--flat-playlist", "--dump-json", "--playlist-end", fmt.Sprintf("%d", paginationThreshold)}
	if browser != "" {
		args = append(args, "--cookies-from-browser", browser)
	}
	args = append(args, feedSubscriptionsURL)

	cmd := exec.Command(YTDLPCommand, args...)
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe for subscriptions: %w", err)
	}

	// Capture stderr to provide better error messages
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stderr pipe for subscriptions: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start yt-dlp command for subscriptions: %w", err)
	}

	var videos []YTVideo
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var currentYTDLPVideo YTDLPVideo
		line := scanner.Bytes()

		if err := json.Unmarshal(line, &currentYTDLPVideo); err != nil {
			// It's possible yt-dlp outputs non-json lines (e.g., warnings about cookies)
			// For now, continue
			continue
		}

		videos = append(videos, newYTVideoFromYTDLPVideo(currentYTDLPVideo))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading yt-dlp subscription output: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		stderrBytes, _ := io.ReadAll(stderr)
		// Provide more informative error message including yt-dlp's output
		return nil, fmt.Errorf("yt-dlp command for subscriptions failed: %w\n%s", err, string(stderrBytes))
	}

	return videos, nil
}