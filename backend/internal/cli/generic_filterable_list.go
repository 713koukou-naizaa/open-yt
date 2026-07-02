package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// A generic, filterable list model
type FilterableListModel struct {
	allItems             []string
	filteredItems        []string
	cursorPosition       int
	selectedChoice       string
	shouldGoBack         bool
	listPrompt           string
	filterQuery          string
	terminalWindowWidth  int
	terminalWindowHeight int
}

func NewFilterableListModel(allItems []string, listPrompt string) FilterableListModel {
	return FilterableListModel{
		allItems:      allItems,
		filteredItems: allItems,
		listPrompt:    listPrompt,
	}
}

func (filterableList FilterableListModel) Init() tea.Cmd {
	return nil
}

func (filterableList FilterableListModel) Update(userMessage tea.Msg) (tea.Model, tea.Cmd) {
	switch userMessageType := userMessage.(type) {
	case tea.WindowSizeMsg:
		filterableList.terminalWindowWidth = userMessageType.Width
		filterableList.terminalWindowHeight = userMessageType.Height
		return filterableList, nil

	case tea.KeyMsg:
		switch userMessageType.String() {
		case "ctrl+c":
			return filterableList, tea.Quit

		case "esc":
			filterableList.shouldGoBack = true
			return filterableList, tea.Quit

		case "enter":
			if len(filterableList.filteredItems) > 0 && filterableList.cursorPosition < len(filterableList.filteredItems) {
				filterableList.selectedChoice = filterableList.filteredItems[filterableList.cursorPosition]
			}
			return filterableList, tea.Quit

		case "backspace":
			if len(filterableList.filterQuery) > 0 {
				filterableList.filterQuery = filterableList.filterQuery[:len(filterableList.filterQuery)-1]
				filterableList.filterItems()
			}

		case "up", "left":
			if len(filterableList.filteredItems) > 0 {
				if filterableList.cursorPosition > 0 {
					filterableList.cursorPosition--
				} else {
					filterableList.cursorPosition = len(filterableList.filteredItems) - 1
				}
			}

		case "down", "right":
			if len(filterableList.filteredItems) > 0 {
				if filterableList.cursorPosition < len(filterableList.filteredItems)-1 {
					filterableList.cursorPosition++
				} else {
					filterableList.cursorPosition = 0
				}
			}

		default:
			if userMessageType.Type == tea.KeyRunes {
				filterableList.filterQuery += string(userMessageType.Runes)
				filterableList.filterItems()
			}
		}
	}
	return filterableList, nil
}

func (filterableList FilterableListModel) View() string {
	if filterableList.terminalWindowWidth == 0 {
		return "Initializing..."
	}

	var filterableListStringView strings.Builder
	filterableListStringView.WriteString(fmt.Sprintf("%s\nFilter: %s█\n\n", filterableList.listPrompt, filterableList.filterQuery))

	// Calculate viewport
	headerHeight := 4 // Lines for prompt, filter, and newlines
	footerHeight := 2 // Lines for help text
	listTotalHeight := filterableList.terminalWindowHeight - headerHeight - footerHeight

	start := 0
	end := len(filterableList.filteredItems)

	if len(filterableList.filteredItems) > listTotalHeight {
		start = filterableList.cursorPosition - listTotalHeight/2
		if start < 0 {
			start = 0
		}
		end = start + listTotalHeight
		if end > len(filterableList.filteredItems) {
			end = len(filterableList.filteredItems)
			start = end - listTotalHeight
		}
	}

	for index := start; index < end; index++ {
		itemString := filterableList.filteredItems[index]
		cursorString := " "
		style := normalItemStyle
		if filterableList.cursorPosition == index {
			cursorString = ">"
			style = selectedItemStyle
		}
		filterableListStringView.WriteString(style.Render(fmt.Sprintf("%s %s", cursorString, itemString)) + "\n")
	}

	if len(filterableList.filteredItems) == 0 {
		filterableListStringView.WriteString("No items match your filter.\n")
	}

	filterableListStringView.WriteString("\n(type to filter, arrows to move, enter to select, esc to go back, CTRL+C to quit)\n")
	return filterableListStringView.String()
}

func (filterableList *FilterableListModel) filterItems() {
	lowerQuery := strings.ToLower(filterableList.filterQuery)
	filteredItems := []string{}
	for _, itemString := range filterableList.allItems {
		if strings.Contains(strings.ToLower(itemString), lowerQuery) {
			filteredItems = append(filteredItems, itemString)
		}
	}
	filterableList.filteredItems = filteredItems
	filterableList.cursorPosition = 0 // Reset cursor
}
