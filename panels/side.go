package panels

import (
	global "htty/globals"
	"htty/panels/components"
	"htty/types"
	utils "htty/utils"

	lipgloss "charm.land/lipgloss/v2"
	tea "github.com/charmbracelet/bubbletea"
)

type SidePane struct {
	fileTree components.FileTree
	
	Dimensions types.PaneGeometry 
	PaneConfig types.HttyPanel 

}

func (side *SidePane) Init() tea.Cmd {
	side.PaneConfig = global.Config.Panels.Side
	side.fileTree = components.FileTree{
		PanelID:    global.PANEL_SIDE_ID,
		PanelTitle: "filetree",
		Border:     types.BorderConfig{Top: true, Bottom: true, Left: true, Right: true},
	}
	return side.fileTree.Init()
}

func (side *SidePane) Update(msg tea.Msg) tea.Cmd {
    focused := global.CurrentPanelID == global.PANEL_FOCUS_IDS[global.PANEL_SIDE_ID]
    if focused {
        utils.SetStatusLineOptions([]string{})
    }
    cmd := side.fileTree.Update(msg)
    if didSelect, path := side.fileTree.Picker.DidSelectFile(msg); didSelect {
        FileTreeHandler(path)
    }
    return cmd
}

func (side SidePane) View() string {
	border := types.BorderConfig{Top: true, Bottom: true, Left: true, Right: true, Width: side.Dimensions.Width - 2}
	focusColor := utils.GetPanelFocusColor(global.PANEL_SIDE_ID)
	style := utils.SetBorder(border).
		BorderForeground(lipgloss.Color(focusColor))
	return utils.SetBorderStyle_WithLabelTop(style, side.fileTree.View(), border,
		utils.GetPanelTitleLabel("filetree", global.PANEL_FOCUS_IDS[global.PANEL_SIDE_ID]),
	)
}

func (side *SidePane) SetSize() {
	side.fileTree.Dimensions.Width = side.Dimensions.Width - 2
	side.fileTree.Dimensions.Height = side.Dimensions.Height - 4
	side.fileTree.SetSize()
}

func FileTreeHandler(path string) {
	global.StateBus.Publish(global.EVENT_STATE_LOAD, utils.LoadState(path))
}
