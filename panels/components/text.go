package components 

import (
	global "htty/globals"
	types "htty/types"
	utils "htty/utils"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type TextPane struct {
	Width, Height int
	Input         textarea.Model
	CharLimit     int
	PanelTitle    string
	PanelID       string
	Placeholder   string
	Showline      bool
	Border        types.BorderConfig
	Margin        types.MarginConfig
	StatusOptions []string

	Dimensions    types.PaneGeometry
	PaneCfg       types.HttyPanel
}

func (text *TextPane) Init() tea.Cmd {
	text.CharLimit = 2052635
	text.Border = types.BorderConfig{Top: true, Bottom: true, Left: true, Right: true}
	text.Showline = true 
	text.StatusOptions = []string{}

	var input textarea.Model = textarea.New()
	input.Placeholder = text.Placeholder
	input.ShowLineNumbers = text.Showline
	input.SetHeight(text.Height)
	input.CharLimit = text.CharLimit
	input.Prompt = ""
	text.Input = input	
	return nil
}

func (text *TextPane) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	focused := global.CurrentPanelID == global.PANEL_FOCUS_IDS[text.PanelID]
	//set a line focus pointer when this panel is focused on
	if focused {
		text.Input.Focus()
		text.Input.Prompt = "│ "
		utils.SetStatusLineOptions(text.StatusOptions)
	} else {
		text.Input.Blur()
		text.Input.Prompt = ""
	}
	text.Input, cmd = text.Input.Update(msg)
	return cmd
}

func (text TextPane) View() string {
	text.Border.Color = utils.GetPanelFocusColor(text.PanelID) 
	style := utils.SetBorder(text.Border).
				BorderForeground(lipgloss.Color(text.Border.Color)).
				Background(lipgloss.Color(global.Config.Common.Background_color))
	return utils.SetBorderStyle_WithLabelTop(style, text.Input.View(), text.Border, 
		utils.GetPanelTitleLabel(text.PanelTitle, global.PANEL_FOCUS_IDS[text.PanelID]),		
	)
}

func (text *TextPane) SetSize() {
	text.Input.SetWidth(text.Dimensions.Width)
	text.Input.SetHeight(text.Dimensions.Height)
}


func (text *TextPane) SetValue(value string){
	text.Input.SetValue(value)
}
