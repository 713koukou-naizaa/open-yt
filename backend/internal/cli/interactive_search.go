package cli

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"open-yt/internal/youtube"
)

// Holds state for interactive search list
type YTSearchModel struct {
	FilterableListModel
	allVideos     []youtube.YTVideo
	videoMap      map[string]youtube.YTVideo
	selectedVideo *youtube.YTVideo // selected video
}

func newYTSearchModel(returnedVideos []youtube.YTVideo) YTSearchModel {
	videoStrings := make([]string, 0, len(returnedVideos))
	videoMap := make(map[string]youtube.YTVideo)
	for i, video := range returnedVideos {
		displayString := formatVideoForList(video, i)
		videoStrings = append(videoStrings, displayString)
		videoMap[displayString] = video
	}

	return YTSearchModel{
		FilterableListModel: NewFilterableListModel(videoStrings, "Select a video to play:"),
		allVideos:           returnedVideos,
		videoMap:            videoMap,
	}
}

func (m YTSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Use the embedded model's Update
	newModel, cmd := m.FilterableListModel.Update(msg)
	m.FilterableListModel = newModel.(FilterableListModel)

	// If a selection was made, map it back to the YTVideo struct
	if m.selected != "" {
		selectedVid := m.videoMap[m.selected]
		m.selectedVideo = &selectedVid
	}

	return m, cmd
}

func formatVideoForList(video youtube.YTVideo, index int) string {
	// Format duration
	durationMinutes := int(video.Duration / 60)
	durationSeconds := int(video.Duration) % 60
	durationStr := fmt.Sprintf("%d:%02d", durationMinutes, durationSeconds)

	// Render video row
	if video.Channel != "" {
		return fmt.Sprintf("%d. [%s] [%s] %s", index, video.Channel, durationStr, video.Title)
	}
	return fmt.Sprintf("%d. [%s] %s", index, durationStr, video.Title)
}