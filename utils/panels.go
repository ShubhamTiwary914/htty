package htty

import (
	types "htty/types"
	global "htty/globals"

	tea "github.com/charmbracelet/bubbletea"
	"charm.land/lipgloss/v2"
)


func GetPercent(percentage int, source int) int{
	return (percentage * source)/100
}

//move focus of cursor onto next item in panels list 
//(present at types/panels -> PANEL_FOCUS_IDS)
func PanelFocusNext(focusID *int){
	(*focusID)++
	if (*focusID >= global.FOCUSABLE_PANELS){
		(*focusID) = 0
	}
}

//move focus of cursor onto previous item in panels list 
//(present at types/panels -> PANEL_FOCUS_IDS)
func PanelFocusPrev(focusID *int){
	(*focusID)--
	if (*focusID < 0){
		(*focusID) = global.FOCUSABLE_PANELS-1
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

func SetBorder(cfg types.BorderConfig) lipgloss.Style {
	style := lipgloss.NewStyle().
		Width(cfg.Width).
		Height(cfg.Height).
		Background(lipgloss.Color(global.Config.Common.Background_color))
	
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
	
	// CHANGED: Use single Border() call with all sides
	style = style.Border(border, top, right, bottom, left)

	if cfg.Color != "" {
		style = style.BorderForeground(lipgloss.Color(cfg.Color))
	}
	return style
}

func SetFullBorder(width, height int, color string) lipgloss.Style {
	return SetBorder(types.BorderConfig{
		Width:  width,
		Height: height,
		Color:  color,
	})
}

func SetBorderOneSide(width, height int, color string, direction string) lipgloss.Style {
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
	case global.BORDER_UP:
		cfg.Top = true
	case global.BORDER_DOWN:
		cfg.Bottom = true
	case global.BORDER_LEFT:
		cfg.Left = true
	case global.BORDER_RIGHT:
		cfg.Right = true
	}

	return SetBorder(cfg)
}

/*
	Creates a new layer of lipgloss compositor 
	(ref - https://github.com/charmbracelet/lipgloss?tab=readme-ov-file#compositing)

	(pane types.BasePanel) - what panel, since that's View() is also invoked

	(paneCfg) - which panel config, since that configs like margin, size is taken

	(offsets) - x, y offsets for margin if needed

	Returns a new lipgloss layer to be Rendered
*/
func CreateNewLayer(pane types.BasePanel, paneCfg types.HttyPanel, offsets ...int) *lipgloss.Layer {
	var xoff int = 0
	var yoff int = 0
	if(len(offsets) >= 2){
		yoff += offsets[1]
	}
	if(len(offsets) >= 1){
		xoff += offsets[0]
	}
	var newlayer *lipgloss.Layer = lipgloss.NewLayer(pane.View()).
					X(paneCfg.Margin[0] + xoff).
					Y(paneCfg.Margin[1] + yoff).
					Z(paneCfg.Layer)
	return newlayer 
}
