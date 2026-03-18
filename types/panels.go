package htty

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//mother of panel i.e: the interface that every panel type should follow
type BasePanel interface {
	Init() tea.Cmd
	Update(tea.Msg) tea.Cmd
	View() string
}

type BorderConfig struct {
	Enabled bool
	Width int
	Height int
	Border  lipgloss.Border // if zero, use NormalBorder
	Color   lipgloss.Color
	Top    bool
	Bottom bool
	Left   bool
	Right  bool
}

type MarginConfig struct {
	Top    int
	Bottom int 
	Left   int  
	Right  int 
}
