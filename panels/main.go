package htty

import (
	utils "htty/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainPane struct {
	width int	
	height int
	margin int
	requestPane RequestPane 
}

func (main *MainPane) Init() tea.Cmd {	
	utils.Infof("main panel initalization")
	return main.requestPane.Init();
}

func (main *MainPane) Update(msg tea.Msg) (tea.Cmd) {
	return utils.UpdatePanels(msg, &main.requestPane)
}

func (main MainPane) View() string {
	style := lipgloss.NewStyle()
	return style.Render(main.requestPane.View())
}

func (main *MainPane) SetSize(w int, h int, m int) {
	main.width = w
	main.height = h
	main.margin = m
	main.requestPane.SetSize(w, h/2)
}
