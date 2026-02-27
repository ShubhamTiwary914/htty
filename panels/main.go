package htty

import (
	panelutil "htty/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainPane struct {
	width int	
	height int
	margin int
}

func (main MainPane) Init() tea.Cmd {
	return nil;
}

func (main MainPane) Update(msg tea.Msg) (MainPane, tea.Cmd) {
	return main, nil;
}

func (main MainPane) View() string {
	style := panelutil.SetBorder(main.width - main.margin, main.height - main.margin, lipgloss.RoundedBorder())
	return style.Render("main panel")
}

func (main *MainPane) SetSize(w int, h int, m int) {
	main.width = w
	main.height = h
	main.margin = m
}
