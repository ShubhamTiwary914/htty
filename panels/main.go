package htty

import (
	"github.com/charmbracelet/bubbletea"
)

type MainPane struct {
}


func (main MainPane) Init() tea.Cmd {
	return nil;
}

func (main MainPane) Update(msg tea.Msg) (MainPane, tea.Cmd) {
	return main, nil;
}

func (m MainPane) View() string {
	return "main panel"
}

