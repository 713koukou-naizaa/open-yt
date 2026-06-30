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
		var YTDLPVid YTDLPVideo
		line := scanner.Bytes()

		// If yt-dlp outputs are non-json lines
		if err := json.Unmarshal(line, &YTDLPVid); err != nil {
			return nil, fmt.Errorf("failed to unmarshal json: %w\nJSON: %s", err, string(line))
		}

		var thumbnails []VideoThumbnail
		for _, t := range YTDLPVid.Thumbnails {
			thumbnails = append(thumbnails, VideoThumbnail{
				URL:    t.URL,
				Height: t.Height,
				Width:  t.Width,
			})
		}

		videos = append(videos, YTVideo{
			Title:       YTDLPVid.Title,
			URL:         YTDLPVid.URL,
			Description: YTDLPVid.Description,
			Duration:    YTDLPVid.Duration,
			Channel:     YTDLPVid.Channel,
			ViewCount:   YTDLPVid.ViewCount,
			Thumbnails:  thumbnails,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading yt-dlp output: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("yt-dlp command failed: %w", err)
	}

	return videos, nil
}
