package youtube

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
)

func Search(query string, maxResults int, YTDLPCommand string) ([]YTVideo, error) {
	searchQuery := fmt.Sprintf("ytsearch%d:%s", maxResults, query)
	cmd := exec.Command(YTDLPCommand, "--flat-playlist", "--dump-json", searchQuery)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start yt-dlp command: %w", err)
	}

	var videos []YTVideo
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var YTDLPVideo YTDLPVideo
		if err := json.Unmarshal(scanner.Bytes(), &YTDLPVideo); err == nil {
			videos = append(videos, newYTVideoFromYTDLPVideo(YTDLPVideo))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading yt-dlp output: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("yt-dlp command failed: %w", err)
	}

	return videos, nil
}

// Converts YTDLPVideo struct to internal YTVideo struct
func newYTVideoFromYTDLPVideo(YTDLPVideo YTDLPVideo) YTVideo {
	var thumbnails []VideoThumbnail
	for _, t := range YTDLPVideo.Thumbnails {
		thumbnails = append(thumbnails, VideoThumbnail{
			URL:    t.URL,
			Height: t.Height,
			Width:  t.Width,
		})
	}

	return YTVideo{
		ID:          YTDLPVideo.ID,
		Title:       YTDLPVideo.Title,
		URL:         YTDLPVideo.URL,
		Description: YTDLPVideo.Description,
		Duration:    YTDLPVideo.Duration,
		Channel:     YTDLPVideo.Channel,
		ViewCount:   YTDLPVideo.ViewCount,
		Thumbnails:  thumbnails,
	}
}
