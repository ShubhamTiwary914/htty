package htty

import (
	types "htty/types"
	global "htty/globals"
	utils "htty/utils"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "charm.land/lipgloss/v2"
	"strings"
)

type TextOptions struct {
	Width, Height int
	Input         textarea.Model
	CharLimit     int
	PanelID       string
	Placeholder   string
	Showline      bool
	Border        types.BorderConfig
	Margin        types.MarginConfig
	OptionItems   []string // Static options to display
}

func (text *TextOptions) Init() tea.Cmd {
	var input textarea.Model = textarea.New()
	input.Placeholder = text.Placeholder
	input.ShowLineNumbers = text.Showline
	input.CharLimit = text.CharLimit
	input.Prompt = ""
	text.Input = input
	return nil
}

func (text *TextOptions) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	focused := global.CurrentPanelID == global.PANEL_FOCUS_IDS[text.PanelID]
	
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

// View renders just the text input (for use when called directly)
func (text TextOptions) View() string {
	baseView, _ := text.ViewWithOptions(false)
	return baseView
}

// ViewWithOptions returns the base view AND optionally the options layer
// If withLayer is true, returns the options layer for compositor
// If withLayer is false, returns nil layer (used when View() is called directly)
func (text TextOptions) ViewWithOptions(withLayer bool) (baseView string, optionsLayer *lipgloss.Layer) {
	// Render the base text input
	inputStyle := utils.SetBorder(text.Border).BorderForeground(
		lipgloss.Color(utils.GetPanelFocusColor(text.PanelID)),
	).Margin(text.Margin.Top, text.Margin.Right, text.Margin.Bottom, text.Margin.Left)
	
	baseView = inputStyle.Render(text.Input.View())
	
	// If not requesting layer, return early
	if !withLayer {
		return baseView, nil
	}
	
	// Check if this panel is focused
	focused := global.CurrentPanelID == global.PANEL_FOCUS_IDS[text.PanelID]
	
	// If not focused or no options, return nil layer
	if !focused || len(text.OptionItems) == 0 {
		return baseView, nil
	}

	// Create options overlay
	optionsStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true, true, true).
		BorderForeground(lipgloss.Color(global.Config.Common.Focus_border_color)).
		Background(lipgloss.Color(global.Config.Common.Background_color)).
		Width(text.Width).
		Padding(0, 1)

	optionsText := strings.Join(text.OptionItems, "\n")
	optionsBox := optionsStyle.Render(optionsText)

	// Return the options layer with Z=2 (above base panels)
	optionsLayer = lipgloss.NewLayer(optionsBox).Z(2)
	
	return baseView, optionsLayer
}



func (text *TextOptions) SetSize(width, height int) {
	text.Width = width
	text.Height = height
	text.Input.SetWidth(width)
	text.Input.SetHeight(height)
}

func (text *TextOptions) SetOptions(options []string) {
	text.OptionItems = options
}

func (text *TextOptions) ClearOptions() {
	text.OptionItems = nil
}
