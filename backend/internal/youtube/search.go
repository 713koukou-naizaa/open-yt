package youtube

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
)

func Search(searchQuery string, paginationThreshold int, YTDLPCommand string) ([]YTVideo, error) {
	YTSearchSearchQuery := fmt.Sprintf("ytsearch%d:%s", paginationThreshold, searchQuery)
	cmd := exec.Command(YTDLPCommand, "--flat-playlist", "--dump-json", YTSearchSearchQuery)

	stdout, cmdStdoutPipeError := cmd.StdoutPipe()
	if cmdStdoutPipeError != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", cmdStdoutPipeError)
	}

	cmdStartError := cmd.Start()
	if cmdStartError != nil {
		return nil, fmt.Errorf("failed to start yt-dlp command: %w", cmdStartError)
	}

	var YTVideos []YTVideo
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var currentYTDLPVideo YTDLPVideo
		if err := json.Unmarshal(scanner.Bytes(), &currentYTDLPVideo); err == nil {
			YTVideos = append(YTVideos, newYTVideoFromYTDLPVideo(currentYTDLPVideo))
		}
	}

	scannerError := scanner.Err()
	if scannerError != nil {
		return nil, fmt.Errorf("error reading yt-dlp output: %w", scannerError)
	}

	cmdWaitError := cmd.Wait()
	if cmdWaitError != nil {
		return nil, fmt.Errorf("yt-dlp command failed: %w", cmdWaitError)
	}

	return YTVideos, nil
}

// Converts YTDLPVideo struct to YTVideo struct
func newYTVideoFromYTDLPVideo(YTDLPVideo YTDLPVideo) YTVideo {
	var YTDLPVideoThumbnails []VideoThumbnail
	for _, currentYTDLPVideoThumbnail := range YTDLPVideo.Thumbnails {
		YTDLPVideoThumbnails = append(YTDLPVideoThumbnails, VideoThumbnail{
			URL:    currentYTDLPVideoThumbnail.URL,
			Height: currentYTDLPVideoThumbnail.Height,
			Width:  currentYTDLPVideoThumbnail.Width,
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
		Thumbnails:  YTDLPVideoThumbnails,
	}
}
