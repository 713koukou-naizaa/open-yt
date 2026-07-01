package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
)

type menuModel struct {
	choices  []string
	cursor   int
	selected string
}

func newMenuModel() menuModel {
	return menuModel{
		choices: []string{"Search", "Play", "Subscriptions feed", "Quit"},
	}
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "up", "left":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "right":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter":
			m.selected = m.choices[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menuModel) View() string {
	s := "What do you want to do?\n"
	
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\n(arrows to move, enter to select, q / esc / CTRL+C to quit)\n"
	return s
}