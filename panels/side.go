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
}

func (side SidePane) Init() tea.Cmd {
	return nil;
} 

func (side *SidePane) Update(msg tea.Msg) (tea.Cmd) {
	return nil;
}

func (side SidePane) View() string {
	style := utils.SetFullBorder(side.width-2, side.height-2, 
		lipgloss.Color(utils.GetPanelFocusColor(types.PANEL_SIDE_ID)),
	)  
	return style.Render("Side Panel") 
}

func (side *SidePane) SetSize(width int, height int) {
	side.width = width
	side.height = height
}
