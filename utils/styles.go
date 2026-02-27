// shared utilities for panels to use (styling)
package htty

import (
	"github.com/charmbracelet/lipgloss"
)


func SetBorder(width, height int, border lipgloss.Border) lipgloss.Style {
	style := lipgloss.NewStyle().
		Width(width).
		Height(height)
	// default border if zero-value passed
	if border == (lipgloss.Border{}) {
		border = lipgloss.NormalBorder()
	}
	return style.Border(border).BorderForeground(lipgloss.Color("240"))
}


