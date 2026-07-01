package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
)

type SimpleMenuModel struct {
	choices  []string
	cursor   int
	selected string
	prompt   string
}

func NewMainMenuModel() SimpleMenuModel {
	return SimpleMenuModel{
		choices: []string{
			menuSearch,
			menuPlay,
			menuSubscriptions,
			menuSubscriptionsFeed,
			menuQuit,
		},
		prompt:  "What do you want to do?",
	}
}

func NewSimpleMenuModel(choices []string, prompt string) SimpleMenuModel {
	return SimpleMenuModel{
		choices: choices,
		prompt:  prompt,
	}
}

func (m SimpleMenuModel) Init() tea.Cmd {
	return nil
}

func (m SimpleMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "up", "left":
			if len(m.choices) > 0 {
				if m.cursor > 0 {
					m.cursor--
				} else {
					m.cursor = len(m.choices) - 1
				}
			}

		case "down", "right":
			if len(m.choices) > 0 {
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
			}

		case "enter":
			m.selected = m.choices[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m SimpleMenuModel) View() string {
	s := m.prompt + "\n\n"
	
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