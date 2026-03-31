package panels 

import (
	global "htty/globals"
	utils "htty/utils"
	tea "github.com/charmbracelet/bubbletea"
)

type SidePane struct {
	width  int
	height int
}

func (side SidePane) Init() tea.Cmd {
	return nil
}

func (side *SidePane) Update(msg tea.Msg) tea.Cmd {
	focused := global.CurrentPanelID == global.PANEL_FOCUS_IDS[global.PANEL_SIDE_ID]
	if focused {
		utils.SetStatusLineOptions([]string{})
	}
	return nil
}

func (side SidePane) View() string {
	style := utils.SetFullBorder(side.width-2, side.height-2,
		utils.GetPanelFocusColor(global.PANEL_SIDE_ID),
	)
	return style.Render("")
}

func (side *SidePane) SetSize(width int, height int) {
	side.width = width
	side.height = height
}
