package htty

import (
	utils "htty/utils"
	types "htty/types"
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

func (side *SidePane) Update(msg tea.Msg) (tea.Cmd) {
	return nil;
}

func (side SidePane) View() string {
	style := utils.SetFullBorder(side.width-side.margin, side.height-side.margin, 
		lipgloss.Color(utils.GetPanelFocusColor(types.PANEL_SIDE_ID)),
	)  
	return style.Render("Side Panel (to be added)")
}

func (side *SidePane) SetSize(w int, h int, m int) {
	side.width = w
	side.height = h
	side.margin = m
}
