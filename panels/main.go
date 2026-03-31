package panels 

import (
	global "htty/globals"
	utils "htty/utils"

	"time"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "charm.land/lipgloss/v2"
)

type tickMsg struct{}
type httpDoneMsg struct {
    output  string
    raw     string
    headers map[string]string
    status  int
}

type MainPane struct {
	width        int
	height       int
	requestPane  RequestPane
	responsePane ResponsePane
	compositor   *lipgloss.Compositor
	dots string 
}

func (main *MainPane) Init() tea.Cmd {
	utils.Infof("main panel initialization")
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
					return httpDoneMsg{output: utils.ResponseParser_main(resp, headers, status, false), raw: string(resp), headers: headers, status: status}
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
	case httpDoneMsg:
		main.responsePane.SetResponse(msg.output, msg.raw, msg.headers, msg.status)
		main.responsePane.loading = false
	}

	return utils.UpdatePanels(msg, &main.requestPane, &main.responsePane)
}

func (main *MainPane) View() string {
	reqLayer := utils.CreateNewLayer(&main.requestPane, global.Config.Panels.Main_req)
	resLayer := utils.CreateNewLayer(&main.responsePane, global.Config.Panels.Main_res,
		0, utils.GetPercent(global.Config.Panels.Main_req.Height, main.height),
	)
	main.compositor = lipgloss.NewCompositor(reqLayer, resLayer)
	mainStyle := lipgloss.NewStyle().Background(lipgloss.Color(global.Config.Common.Background_color))
	return mainStyle.Render(main.compositor.Render())
}

func (main *MainPane) SetSize(width int, height int) {
	main.width = width
	main.height = height
	main.requestPane.SetSize(
		utils.GetPercent(global.Config.Panels.Main_req.Width, main.width),
		utils.GetPercent(global.Config.Panels.Main_req.Height, main.height),
	)
	main.responsePane.SetSize(
		utils.GetPercent(global.Config.Panels.Main_res.Width, main.width),
		utils.GetPercent(global.Config.Panels.Main_res.Height, main.height),
	)
}
