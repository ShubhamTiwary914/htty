package panels

import (
	global "htty/globals"
	"htty/types"

	utils "htty/utils"
	"time"

	lipgloss "charm.land/lipgloss/v2"
	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg struct{}


type MainPane struct {
	requestPane  RequestPane
	responsePane ResponsePane
	compositor   *lipgloss.Compositor
	dots string 
	
	Dimensions types.PaneGeometry 
	PaneConfig types.HttyPanel 
}

func (main *MainPane) Init() tea.Cmd {
	utils.Infof("main panel initialization")	
	main.PaneConfig = global.Config.Panels.Main
	main.requestPane.Init()
	main.responsePane.Init()
	return nil
}

func (main *MainPane) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == global.Config.Key.Sendapicall {
			main.responsePane.loading = true
			main.dots = ""
			return tea.Batch(
				tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg{} }),
				func() tea.Msg {
					resp, headers, status, err := utils.HTTPCaller(main.requestPane.ExportPayload())
					if err != nil {
						utils.Errorf("error loading response, error: %v", err)
					}
					return types.HttpRespState{
						Output: utils.ResponseParser_main(resp, headers, status, false), 
						Raw: string(resp), Headers: headers, Status: status,
					}
				},
			)
		}

	case tickMsg:
		if main.responsePane.loading {
			main.dots += "."
			main.responsePane.SetResponse("Loading"+ main.dots, "", nil, 0)
			return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg{} })
		}
		return nil;
	case types.HttpRespState:
		main.responsePane.SetResponse(msg.Output, msg.Raw, msg.Headers, msg.Status)
		main.responsePane.loading = false
	}

	return utils.UpdatePanels(msg, &main.requestPane, &main.responsePane)
}

func (main *MainPane) View() string {
	reqLayer := utils.CreateNewLayer(&main.requestPane, global.Config.Panels.Main_req)
	resLayer := utils.CreateNewLayer(&main.responsePane, global.Config.Panels.Main_res,
		0, utils.GetPercent(global.Config.Panels.Main_req.Height, main.Dimensions.Height),
	)
	main.compositor = lipgloss.NewCompositor(reqLayer, resLayer)
	mainStyle := lipgloss.NewStyle().Background(lipgloss.Color(global.Config.Common.Background_color))
	return mainStyle.Render(main.compositor.Render())
}


func (main *MainPane) SetSize() {
	main.requestPane.Dimensions = utils.GetPaneGeometry(main.requestPane.PaneConfig, main.Dimensions)
	main.responsePane.Dimensions = utils.GetPaneGeometry(main.responsePane.PaneConfig, main.Dimensions)
	main.requestPane.SetSize()
	main.responsePane.SetSize()
}
