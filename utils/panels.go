package htty

import (
	types "htty/types"
	global "htty/globals"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


func GetPercent(percentage int, source int) int{
	return (percentage * source)/100
}

//move focus of cursor onto next item in panels list 
//(present at types/panels -> PANEL_FOCUS_IDS)
func PanelFocusNext(focusID *int){
	(*focusID)++
	if (*focusID >= len(global.PANEL_FOCUS_IDS)){
		(*focusID) = 0
	}
}

//move focus of cursor onto previous item in panels list 
//(present at types/panels -> PANEL_FOCUS_IDS)
func PanelFocusPrev(focusID *int){
	(*focusID)--
	if (*focusID < 0){
		(*focusID) = len(global.PANEL_FOCUS_IDS)-1
	}
}


//jump to panel using the panelID(2) or panel mapping string(ex: "PANEL_REQ_METHOD_ID") 
//(present at types/panels -> PANEL_FOCUS_IDS)
//NOTE: using the panel enums constant is recommended, instead of using hardcoded magic numbers
func PanelFocusJump(focusID *int, newfocuskey interface{}){
	switch newfocuskey.(type) {
		case int:
			*focusID = newfocuskey.(int)
		case string:
			*focusID = (global.PANEL_FOCUS_IDS[(newfocuskey.(string))])
	}	
}

/*
	passes tea Msg object to child panels for handling update events, 
	in order words, it allows these child panels to be "updated" to new events
	
	usage:  UpdatePanels(msg, &childpane1, &childpane2...)

	"msg" is the object of type tea.Msg which can hold event actions like keyboard press, mouse clicks etc
*/
func UpdatePanels(msg tea.Msg, panels ...types.BasePanel) tea.Cmd {
	var cmds []tea.Cmd
	for _, p := range panels {
		if cmd := p.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}


// Get focus color(bright) if this panel is focused on, or get unfocus color(dim) if this panel isnt in focus.  
// Can be used for highlighting border of focused panels
// (can pass panelID like costant PANEL_REQ_METHOD_ID present or string like "SIDE")
func GetPanelFocusColor(panelkey interface{}) string {
	var panelID int
	switch panelkey.(type){
		case int:
			panelID = panelkey.(int)
		case string:
			panelID = (global.PANEL_FOCUS_IDS[(panelkey.(string))])
	}
	if global.CurrentPanelID == panelID {
		return global.Config.Common.Focus_border_color 
	}
	return global.Config.Common.Unfocus_border_color
}


// General main method to set lipgloss border for panels
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
