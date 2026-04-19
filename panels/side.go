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
		PanelTitle: "File tree sidebar",
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
	border := types.BorderConfig{Top: true, Bottom: true, Left: true, Right: true}
	focusColor := utils.GetPanelFocusColor(global.PANEL_SIDE_ID)
	style := utils.SetBorder(border).
		BorderForeground(lipgloss.Color(focusColor)).
		Background(lipgloss.Color(global.Config.Common.Background_color))
	return utils.SetBorderStyle_WithLabelTop(style, side.fileTree.View(), border,
		utils.GetPanelTitleLabel("File tree sidebar", global.PANEL_FOCUS_IDS[global.PANEL_SIDE_ID]),
	)
}

func (side *SidePane) SetSize() {
	side.fileTree.SetSize(side.Dimensions.Width, side.Dimensions.Height)
}

func FileTreeHandler(path string) {
	utils.Debugf("FileTreeHandler called - selected path: %s", path)
}
