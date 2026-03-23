package htty

import (
	global "htty/globals"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/bubbletea"
)

type StatusLinePane struct {
	options []string
}


func (status *StatusLinePane) Init() tea.Cmd {
	return nil
}

func (status *StatusLinePane) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (status *StatusLinePane) View() string {
	status.options = global.StatusLineOptions
	textStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color(global.Config.Common.Statusline_options_color)).
					Bold(true)
	var renderedOptions []string
	for _, opt := range status.options {
		renderedOptions = append(renderedOptions, textStyle.Render(opt))
	}
	optionsStr := strings.Join(renderedOptions, "  ")	
	return textStyle.Render(optionsStr)
}
