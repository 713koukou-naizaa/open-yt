package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"open-yt/internal/youtube"
)

// Holds state for interactive search list
type YTSearchModel struct {
	returnedVideos   []youtube.YTVideo // list of returned videos
	cursorPosition   int  // video at which cursor is pointing at
	selectedVideo *youtube.YTVideo // selected video
	back          bool             // if user wants to go back
}

func newYTSearchModel(returnedVideos []youtube.YTVideo) YTSearchModel {
	return YTSearchModel{returnedVideos: returnedVideos}
}

// First function called
// Usually returns optional initial command
// We don't need it, so returning nil
func (m YTSearchModel) Init() tea.Cmd {
	return nil
}

// Called when message received
// Where we'll handle user input and other events
func (m YTSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Message is a key press
	case tea.KeyMsg:

		// Which key
		switch msg.String() {
			// Exit program
			case "ctrl+c":
				return m, tea.Quit

			// Go back to previous menu
			case "q", "b":
				m.back = true
				return m, tea.Quit

			// Move cursor up
			case "up", "left", "k":
				if m.cursorPosition > 0 {
					m.cursorPosition--
				}

			// Move cursor down
			case "down", "right", "j":
				if m.cursorPosition < len(m.returnedVideos)-1 {
					m.cursorPosition++
				}

			// Set selected video and exit
			case "enter":
				m.selectedVideo = &m.returnedVideos[m.cursorPosition]
				return m, tea.Quit
			}
	}

	// Return updated model to Bubble Tea runtime
	return m, nil
}

// View renders UI
// Called every time model is updated
func (m YTSearchModel) View() string {
	s := "Select a video to play:\n"

	// For each returned video
	for i, video := range m.returnedVideos {
		cursor := " " // Hide cursor
		if m.cursorPosition == i {
			cursor = ">" // Display cursor to select visual
		}

		// Format duration
		durationMinutes := int(video.Duration / 60)
		durationSeconds := int(video.Duration) % 60
		durationStr := fmt.Sprintf("%d:%02d", durationMinutes, durationSeconds)

		// Render video row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, durationStr, video.Title)
	}

	s += "\n(arrows or j/k to move, enter to select, q/b to go back, CTRL+C to quit)\n"

	return s
}