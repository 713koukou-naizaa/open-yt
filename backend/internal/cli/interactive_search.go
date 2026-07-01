package cli

import (
	"fmt"

	"strings"
	"github.com/charmbracelet/bubbletea"
	"open-yt/internal/youtube"
)

// Holds state for interactive search list
type YTSearchModel struct {
	returnedVideos []youtube.YTVideo // list of returned videos
	filteredVideos []youtube.YTVideo // filtered videos based on query
	cursorPosition int  // video at which cursor is pointing at
	selectedVideo *youtube.YTVideo // selected video
	back bool // if user wants to go back
	filterQuery string
}

func newYTSearchModel(returnedVideos []youtube.YTVideo) YTSearchModel {
	return YTSearchModel{
		returnedVideos: returnedVideos,
		filteredVideos: returnedVideos, // Initially, all videos are shown
	}
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// Message is a key press
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			if len(m.filteredVideos) > 0 && m.cursorPosition < len(m.filteredVideos) {
				m.selectedVideo = &m.filteredVideos[m.cursorPosition]
			}
			return m, tea.Quit

		case "backspace":
			if len(m.filterQuery) > 0 {
				m.filterQuery = m.filterQuery[:len(m.filterQuery)-1]
				m.filterVideos()
			}

		case "esc":
			m.back = true
			return m, tea.Quit
		case "up", "left":
			if m.cursorPosition > 0 {
				m.cursorPosition--
			}

		case "down", "right":
			if m.cursorPosition < len(m.filteredVideos)-1 {
				m.cursorPosition++
			}
		default:
			// Handle typing for filter
			if msg.Type == tea.KeyRunes {
				m.filterQuery += string(msg.Runes)
				m.filterVideos()
			}
		}
	}

	// Return updated model to Bubble Tea runtime
	return m, cmd
}

// Updates filteredVideos slice based on filterQuery
func (m *YTSearchModel) filterVideos() {
	lowerQuery := strings.ToLower(m.filterQuery)
	filtered := []youtube.YTVideo{}
	for _, video := range m.returnedVideos {
		if strings.Contains(strings.ToLower(video.Title), lowerQuery) || strings.Contains(strings.ToLower(video.Channel), lowerQuery) {
			filtered = append(filtered, video)
		}
	}
	m.filteredVideos = filtered
	m.cursorPosition = 0 // Reset cursor
}

// View renders UI
// Called every time model is updated
func (m YTSearchModel) View() string {
	s := fmt.Sprintf("Select a video to play.\nFilter: %s█\n\n", m.filterQuery)

	// For each returned video
	for i, video := range m.filteredVideos {
		cursor := " " // Hide cursor
		if m.cursorPosition == i {
			cursor = ">" // Display cursor to select visual
		}

		// Format duration
		durationMinutes := int(video.Duration / 60)
		durationSeconds := int(video.Duration) % 60
		durationStr := fmt.Sprintf("%d:%02d", durationMinutes, durationSeconds)

		// Render video row
		if video.Channel != "" {
			s += fmt.Sprintf("%s %d. [%s] [%s] %s\n", cursor, i, video.Channel, durationStr, video.Title)
		} else {
			s += fmt.Sprintf("%s %d. [%s] %s\n", cursor, i, durationStr, video.Title)
		}
	}

	if len(m.filteredVideos) == 0 {
		s += "No videos match your filter.\n"
	}

	s += "\n(type to filter, arrows to move, enter to select, esc to go back, CTRL+C to quit)\n"

	return s
}