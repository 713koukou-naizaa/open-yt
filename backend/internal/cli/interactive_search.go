package cli

import (
	"fmt"
	"open-yt/internal/youtube"

	tea "github.com/charmbracelet/bubbletea"
)

// Holds state for interactive search list
type YTSearchModel struct {
	FilterableListModel
	allYTVideos           []youtube.YTVideo
	videoStringYTVideoMap map[string]youtube.YTVideo
	selectedYTVideo       *youtube.YTVideo
}

func newYTSearchModel(returnedYTVideos []youtube.YTVideo) YTSearchModel {
	videoStrings := make([]string, 0, len(returnedYTVideos))
	videoStringYTVideoMap := make(map[string]youtube.YTVideo)
	for videoIndex, YTVideo := range returnedYTVideos {
		displayString := formatVideoForList(YTVideo, videoIndex)
		videoStrings = append(videoStrings, displayString)
		videoStringYTVideoMap[displayString] = YTVideo
	}

	return YTSearchModel{
		FilterableListModel:   NewFilterableListModel(videoStrings, "Select a video to play:"),
		allYTVideos:           returnedYTVideos,
		videoStringYTVideoMap: videoStringYTVideoMap,
	}
}

func (YTSearchFilterableList YTSearchModel) Update(userMessage tea.Msg) (tea.Model, tea.Cmd) {
	// Use the embedded model's Update
	updatedYTSearchFilterableList, cmd := YTSearchFilterableList.FilterableListModel.Update(userMessage)
	YTSearchFilterableList.FilterableListModel = updatedYTSearchFilterableList.(FilterableListModel)

	// If a selection was made, map it back to the YTVideo struct
	if YTSearchFilterableList.selectedChoice != "" {
		selectedYTVideo := YTSearchFilterableList.videoStringYTVideoMap[YTSearchFilterableList.selectedChoice]
		YTSearchFilterableList.selectedYTVideo = &selectedYTVideo
	}

	return YTSearchFilterableList, cmd
}

func formatVideoForList(YTVideo youtube.YTVideo, videoIndex int) string {
	// Format duration
	durationMinutes := int(YTVideo.Duration / 60)
	durationSeconds := int(YTVideo.Duration) % 60
	durationString := fmt.Sprintf("%d:%02d", durationMinutes, durationSeconds)

	// Render video row
	if YTVideo.Channel != "" {
		return fmt.Sprintf("%d. [%s] [%s] %s", videoIndex, YTVideo.Channel, durationString, YTVideo.Title)
	}
	return fmt.Sprintf("%d. [%s] %s", videoIndex, durationString, YTVideo.Title)
}
