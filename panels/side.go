package htty

import (
	panelutil "htty/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SidePane struct {
	width int 
	height int
	margin int
}

func (side SidePane) Init() tea.Cmd {
	return nil;
} 

func (side SidePane) Update(msg tea.Msg) (SidePane, tea.Cmd) {
	return side, nil;
}

func (side SidePane) View() string {
	style := panelutil.SetBorder(side.width-side.margin, side.height-side.margin, lipgloss.RoundedBorder())
	return style.Render("side panel")
}

func (side *SidePane) SetSize(w int, h int, m int) {
	side.width = w
	side.height = h
	side.margin = m
}
