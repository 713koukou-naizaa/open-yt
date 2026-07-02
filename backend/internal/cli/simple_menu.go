package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SimpleMenuModel struct {
	menuChoices    []string
	cursorPosition int
	selectedChoice string
	menuPrompt     string
}

var (
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Background(lipgloss.Color("236"))
	normalItemStyle   = lipgloss.NewStyle()
)

func NewMainMenuModel() SimpleMenuModel {
	return SimpleMenuModel{
		menuChoices: []string{
			menuHomeFeed,
			menuSubscriptions,
			menuSubscriptionsFeed,
			menuSearch,
			menuPlay,
			menuQuit,
		},
		menuPrompt: "What do you want to do?",
	}
}

func NewSimpleMenuModel(menuChoices []string, menuPrompt string) SimpleMenuModel {
	return SimpleMenuModel{
		menuChoices: menuChoices,
		menuPrompt:  menuPrompt,
	}
}

func (simpleMenu SimpleMenuModel) Init() tea.Cmd {
	return nil
}

func (simpleMenu SimpleMenuModel) Update(userMessage tea.Msg) (tea.Model, tea.Cmd) {
	switch userMessageType := userMessage.(type) {
	case tea.KeyMsg:
		switch userMessageType.String() {
		case "ctrl+c", "q", "esc":
			return simpleMenu, tea.Quit

		case "up", "left":
			if len(simpleMenu.menuChoices) > 0 {
				if simpleMenu.cursorPosition > 0 {
					simpleMenu.cursorPosition--
				} else {
					simpleMenu.cursorPosition = len(simpleMenu.menuChoices) - 1
				}
			}

		case "down", "right":
			if len(simpleMenu.menuChoices) > 0 {
				if simpleMenu.cursorPosition < len(simpleMenu.menuChoices)-1 {
					simpleMenu.cursorPosition++
				} else {
					simpleMenu.cursorPosition = 0
				}
			}

		case "enter":
			simpleMenu.selectedChoice = simpleMenu.menuChoices[simpleMenu.cursorPosition]
			return simpleMenu, tea.Quit
		}
	}
	return simpleMenu, nil
}

func (simpleMenu SimpleMenuModel) View() string {
	menuStringView := simpleMenu.menuPrompt + "\n\n"

	for index, menuChoice := range simpleMenu.menuChoices {
		menuCursor := " "
		style := normalItemStyle
		if simpleMenu.cursorPosition == index {
			menuCursor = ">"
			style = selectedItemStyle
		}
		menuStringView += style.Render(fmt.Sprintf("%s %s", menuCursor, menuChoice)) + "\n"
	}
	menuStringView += "\n(arrows to move, enter to select, q / esc / CTRL+C to quit)\n"
	return menuStringView
}
