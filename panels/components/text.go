package htty

import (
	types "htty/types"
	utils "htty/utils"

	global "htty/globals"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)
var options = []string{"foo", "bar", "baz"}

type TextModel struct {
	Width, Height int
	Input textarea.Model
	CharLimit int
	PanelID string
	Placeholder string 
	Showline bool 
	Border types.BorderConfig
}

func (text *TextModel) Init() tea.Cmd {
	var input textarea.Model = textarea.New()
	input.Placeholder = text.Placeholder
	input.ShowLineNumbers = text.Showline
	input.SetHeight(text.Height)
	input.CharLimit = text.CharLimit 
	input.Prompt = ""
	text.Input= input	
	return nil
}

func (text *TextModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	focused := global.CurrentPanelID == types.PANEL_FOCUS_IDS[text.PanelID]
	if focused {
		text.Input.Focus()
		text.Input.Prompt = "│ " 
	} else {
		text.Input.Blur()
		text.Input.Prompt = ""
	}
	text.Input, cmd = text.Input.Update(msg)
	return cmd
}

func (text TextModel) View() string {	
	utils.Debugf("borders: (top:%t)(bot:%t)(left:%t)(right:%t)", 
		text.Border.Top,
		text.Border.Bottom,
		text.Border.Left,
		text.Border.Right,
	)
	style := utils.SetBorder(text.Border)
	return style.Render(text.Input.View())
}

func (text *TextModel) SetSize(width, height int) {
	text.Width = width
	text.Height = height
	text.Input.SetWidth(width)
	text.Input.SetHeight(height)
}
