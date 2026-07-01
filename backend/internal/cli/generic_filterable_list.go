package cli

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"strings"
)

// A generic, filterable list model.
type FilterableListModel struct {
	allItems      []string
	filteredItems []string
	cursor        int
	selected      string
	back          bool
	prompt        string
	filterQuery   string
	width         int
	height        int
}

func NewFilterableListModel(items []string, prompt string) FilterableListModel {
	return FilterableListModel{
		allItems:      items,
		filteredItems: items,
		prompt:        prompt,
	}
}

func (m FilterableListModel) Init() tea.Cmd {
	return nil
}

func (m FilterableListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.back = true
			return m, tea.Quit

		case "enter":
			if len(m.filteredItems) > 0 && m.cursor < len(m.filteredItems) {
				m.selected = m.filteredItems[m.cursor]
			}
			return m, tea.Quit

		case "backspace":
			if len(m.filterQuery) > 0 {
				m.filterQuery = m.filterQuery[:len(m.filterQuery)-1]
				m.filterItems()
			}

		case "up", "left":
			if len(m.filteredItems) > 0 {
				if m.cursor > 0 {
					m.cursor--
				} else {
					m.cursor = len(m.filteredItems) - 1
				}
			}

		case "down", "right":
			if len(m.filteredItems) > 0 {
				if m.cursor < len(m.filteredItems)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
			}

		default:
			if msg.Type == tea.KeyRunes {
				m.filterQuery += string(msg.Runes)
				m.filterItems()
			}
		}
	}
	return m, nil
}

func (m FilterableListModel) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	var s strings.Builder
	s.WriteString(fmt.Sprintf("%s\nFilter: %s█\n\n", m.prompt, m.filterQuery))

	// Calculate viewport
	headerHeight := 4 // Lines for prompt, filter, and newlines
	footerHeight := 2 // Lines for help text
	listHeight := m.height - headerHeight - footerHeight

	start := 0
	end := len(m.filteredItems)

	if len(m.filteredItems) > listHeight {
		start = m.cursor - listHeight/2
		if start < 0 {
			start = 0
		}
		end = start + listHeight
		if end > len(m.filteredItems) {
			end = len(m.filteredItems)
			start = end - listHeight
		}
	}

	for i := start; i < end; i++ {
		item := m.filteredItems[i]
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, item))
	}

	if len(m.filteredItems) == 0 {
		s.WriteString("No items match your filter.\n")
	}

	s.WriteString("\n(type to filter, arrows to move, enter to select, esc to go back, CTRL+C to quit)\n")
	return s.String()
}

func (m *FilterableListModel) filterItems() {
	lowerQuery := strings.ToLower(m.filterQuery)
	filtered := []string{}
	for _, item := range m.allItems {
		if strings.Contains(strings.ToLower(item), lowerQuery) {
			filtered = append(filtered, item)
		}
	}
	m.filteredItems = filtered
	m.cursor = 0 // Reset cursor
}