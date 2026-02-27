package htty

import (
	types "htty/types"

	global "htty/globals"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


func PanelFocusNext(focusID *int){
	(*focusID)++
	if (*focusID >= len(types.PANEL_FOCUS_IDS)){
		(*focusID) = 0
	}
}
func PanelFocusPrev(focusID *int){
	(*focusID)--
	if (*focusID < 0){
		(*focusID) = len(types.PANEL_FOCUS_IDS)-1
	}
}

//jump to panel using the panelID or panel mapping string in types/panels/PANEL_FOCUS
func PanelFocusJump(focusID *int, newfocuskey interface{}){
	switch newfocuskey.(type) {
		case int:
			*focusID = newfocuskey.(int)
		case string:
			*focusID = (types.PANEL_FOCUS_IDS[(newfocuskey.(string))])
	}	
}

func UpdatePanels(msg tea.Msg, panels ...types.BasePanel) tea.Cmd {
	var cmds []tea.Cmd
	for _, p := range panels {
		if cmd := p.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

//get color for panel border depdening on whether it is in focus.
//(can pass panelID like 0,1 or string like "SIDE")
func GetPanelFocusColor(panelkey interface{}) string {
	var panelID int
	switch panelkey.(type){
		case int:
			panelID = panelkey.(int)
		case string:
			panelID = (types.PANEL_FOCUS_IDS[(panelkey.(string))])
	}
	if global.CurrentPanelID == panelID {
		return types.PANEL_FOCUS_COLOR
	}
	return types.PANEL_UNFOCUS_COLOR
}


func SetBorder(cfg types.BorderConfig) lipgloss.Style {
	style := lipgloss.NewStyle().
		Width(cfg.Width).
		Height(cfg.Height)
	// default enabled = true
	enabled := true
	if cfg.Enabled {
		enabled = true
	}
	if !enabled {
		return style
	}
	border := cfg.Border
	if border == (lipgloss.Border{}) {
		border = lipgloss.NormalBorder()
	}
	// default sides = true unless explicitly set
	top := true
	bottom := true
	left := true
	right := true
	if cfg.Top || cfg.Bottom || cfg.Left || cfg.Right {
		top = cfg.Top
		bottom = cfg.Bottom
		left = cfg.Left
		right = cfg.Right
	}
	style = style.
		Border(border).
		BorderTop(top).
		BorderBottom(bottom).
		BorderLeft(left).
		BorderRight(right)

	if cfg.Color != "" {
		style = style.BorderForeground(cfg.Color)
	}
	return style
}


func SetFullBorder(width, height int, color lipgloss.Color) lipgloss.Style {
	return SetBorder(types.BorderConfig{
		Width:  width,
		Height: height,
		Color:  color,
	})
}

func SetBorderOneSide(width, height int, color lipgloss.Color, direction string) lipgloss.Style {
	cfg := types.BorderConfig{
		Width:  width,
		Height: height,
		Color:  color,
		Top:    false,
		Bottom: false,
		Left:   false,
		Right:  false,
	}
	switch direction {
	case types.BORDER_UP:
		cfg.Top = true
	case types.BORDER_DOWN:
		cfg.Bottom = true
	case types.BORDER_LEFT:
		cfg.Left = true
	case types.BORDER_RIGHT:
		cfg.Right = true
	}

	return SetBorder(cfg)
}
