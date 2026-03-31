package utils 

import (
	"fmt"
	global "htty/globals"
	types "htty/types"

	"strconv"
	"strings"
	"regexp"

	"charm.land/lipgloss/v2"
	tea "github.com/charmbracelet/bubbletea"
)

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

func GetPanelTitleLabel(title string, id int) string {
	return fmt.Sprintf(" [%d] %s ", id, title)
}

func SetStatusLineOptions(options []string){
	global.StatusLineOptions = options
}


//for a panel, what are the action keys , like alt+s, alt+c allowed as per its "Keys"
func GetPanelActionKeys(panelCfg types.HttyPanel) []string{
	var keys []string;
	for _, key := range panelCfg.Keys {
		keys = append(keys, key)
	}
	return keys 
}


//from panel's key (actions that appear on status line), generates string arr of options to show in statusline
func GetPanel_KeyOptions(panelcfg types.HttyPanel) []string{
	var options []string
	for action, key := range panelcfg.Keys {
		options = append(options, fmt.Sprintf("%s(%s)", action, key))
	}
	return options
}

//given a msg.Sting() from bubble tea.cmd, sees if event is a jump type 
//jump type means tp jump to a panel with <key>+<number> (ex: alt+2)
func EventIs_TypeJumpPanel(eventstr string) (bool, int) {
	re := regexp.MustCompile(`^`+ global.Config.Key.Jumpleader +`\+([0-9])$`)
	if matches := re.FindStringSubmatch(eventstr); matches != nil {
		panelNum, _ := strconv.Atoi(matches[1])
		return true, panelNum
	}
	return false, 0
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


func SetBorderStyle_WithLabel(style lipgloss.Style, content string, cfg types.BorderConfig, label string, position string) string {
    if label == "" {
        return style.Render(content)
    }
    
    isTop := position == "top"
    isBottom := position == "bottom"
    
    if (isTop && !cfg.Top) || (isBottom && !cfg.Bottom) {
        return style.Render(content)
    }
    
    // render content with all borders EXCEPT the one we're labeling
    cfgModified := cfg
    if isTop {
        cfgModified.Top = false
    } else {
        cfgModified.Bottom = false
    }
    
    // get styling from the original style
    borderFgColor := style.GetBorderTopForeground()
    contentWithBorders := SetBorder(cfgModified).
        BorderForeground(borderFgColor).
        Margin(style.GetMarginTop(), style.GetMarginRight(), style.GetMarginBottom(), style.GetMarginLeft()).
        Background(style.GetBackground()).
        Render(content)
    
    // manually construct labeled border
    border := cfg.Border
    if border == (lipgloss.Border{}) {
        border = lipgloss.NormalBorder()
    }
    color := lipgloss.Color(cfg.Color)
    
    var leftCorner, rightCorner, horizontal string
    if isTop {
        leftCorner = border.TopLeft
        rightCorner = border.TopRight
        horizontal = border.Top
    } else {
        leftCorner = border.BottomLeft
        rightCorner = border.BottomRight
        horizontal = border.Bottom
    }
    
    // calculate width from the rendered content
    contentWidth := lipgloss.Width(contentWithBorders)
    remainingWidth := max(contentWidth - lipgloss.Width(label) - 2, 0)
    
    // build the entire line FIRST, then style it all at once
    borderLine := leftCorner + label + strings.Repeat(horizontal, remainingWidth) + rightCorner
    labeledBorder := lipgloss.NewStyle().Foreground(color).Render(borderLine)
    
    var result string
    if isTop {
        result = lipgloss.JoinVertical(lipgloss.Left, labeledBorder, contentWithBorders)
    } else {
        result = lipgloss.JoinVertical(lipgloss.Left, contentWithBorders, labeledBorder)
    }
    return result
}


func SetBorderStyle_WithLabelTop(style lipgloss.Style, content string, cfg types.BorderConfig, label string) string {
    return SetBorderStyle_WithLabel(style, content, cfg, label, "top")
}

func SetBorderStyle_WithLabelBottom(style lipgloss.Style, content string, cfg types.BorderConfig, label string) string {
    return SetBorderStyle_WithLabel(style, content, cfg, label, "bottom")
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
