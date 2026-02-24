package htty

import tea "github.com/charmbracelet/bubbletea"


type SidePane struct {
}

func (side SidePane) Init() tea.Cmd {
	return nil;
} 

func (side SidePane) Update(msg tea.Msg) (SidePane, tea.Cmd) {
	return side, nil;
}

func (side SidePane) View() string {
	return "side panel"
}


