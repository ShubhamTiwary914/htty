package htty

import (
	"fmt"
	utils "htty/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainPane struct {
	width int	
	height int
	requestPane RequestPane 
	responsePane ResponsePane
}

func (main *MainPane) Init() tea.Cmd {	
	utils.Infof("main panel initalization")
	main.requestPane.Init();
	main.responsePane.Init()
	return nil;
}

func (main *MainPane) Update(msg tea.Msg) (tea.Cmd) {
	switch msg := msg.(type) {
	//auto resize to window dimensions
	case tea.KeyMsg:
		//take out inputs from reqPane -> http call -> response to responsePane
		if msg.String() == "alt+enter" {
			resp, status, err := utils.HTTPCaller(main.requestPane.ExportPayload())
			var output string
			if err != nil {
				output = fmt.Sprintf("status: %d\nerror: %v\nresponse:\n", status, err)
			} else {
				output = fmt.Sprintf("status: %d\nerror: nil\nresponse:\n%s", status, string(resp))
			}; 
			main.responsePane.SetResponse(output)
			utils.Debugf(output)
		}
	}	
	return utils.UpdatePanels(msg, &main.requestPane, &main.responsePane)
}

func (main MainPane) View() string {
	style := lipgloss.NewStyle().Margin()
	mainSubPanels := lipgloss.JoinVertical(
		lipgloss.Left,
		main.requestPane.View(),
		main.responsePane.View(),
	)
	return style.Render(mainSubPanels)
}

func (main *MainPane) SetSize(w int, h int, m int) {
	main.width = w
	main.height = h
	main.requestPane.SetSize(w, h/2)
	main.responsePane.SetSize(w, int(float64(h)/2.5))
}

