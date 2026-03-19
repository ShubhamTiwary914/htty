package htty

import (
	types "htty/types"

	global "htty/globals"
	utils "htty/utils"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "charm.land/lipgloss/v2"
)

type TextPane struct {
	Width, Height int
	Input         textarea.Model
	CharLimit     int
	PanelID       string
	Placeholder   string
	Showline      bool
	Border        types.BorderConfig
	Margin        types.MarginConfig
}

func (text *TextPane) Init() tea.Cmd {
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
	} else {
		text.Input.Blur()
		text.Input.Prompt = ""
	}
	text.Input, cmd = text.Input.Update(msg)
	return cmd
}

func (text TextPane) View() string {
	style := utils.SetBorder(text.Border).BorderForeground(
		//adds highlight on border when focused on
		lipgloss.Color(utils.GetPanelFocusColor(text.PanelID)),
	).Margin(
		text.Margin.Top, text.Margin.Right, text.Margin.Bottom, text.Margin.Left)
	return style.Render(text.Input.View())
}

func (text *TextPane) SetSize(width, height int) {
	text.Width = width
	text.Height = height
	text.Input.SetWidth(width)
	text.Input.SetHeight(height)
}
