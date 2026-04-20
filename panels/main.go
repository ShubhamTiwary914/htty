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
	stateLoadChan chan any 
	currStateBuf types.HttyState
}

func (main *MainPane) Init() tea.Cmd {
	utils.Infof("main panel initialization")	
	main.PaneConfig = global.Config.Panels.Main
	main.requestPane.Init()
	main.responsePane.Init()

	main.stateLoadChan = make(chan any, 1)
	global.StateBus.Subscribe(global.EVENT_STATE_LOAD, main.stateLoadChan)
	return main.__waitForStateLoad()
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
					main.currStateBuf.Response = types.HttpRespState{
						Output: utils.ResponseParser_main(resp, headers, status, false), 
						Raw: string(resp), Headers: headers, Status: status,
					}
					return main.currStateBuf.Response
				},
			)
		}
		utils.Debugf("message: %s", msg.String())
		if msg.String() == global.Config.Key.Savestate {
			main.__saveStateDialog()			
			return nil;
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
	case types.HttyState:
		main.__stateLoadPerform(msg)
	}
	return utils.UpdatePanels(msg, &main.requestPane, &main.responsePane)
}


func (main *MainPane) View() string {
	reqLayer := utils.CreateLayerFromDims(&main.requestPane, main.requestPane.Dimensions, 1)
	resLayer := utils.CreateLayerFromDims(&main.responsePane, main.responsePane.Dimensions, 1)

	main.compositor = lipgloss.NewCompositor(reqLayer, resLayer)
	mainStyle := lipgloss.NewStyle().Background(lipgloss.Color(global.Config.Common.Background_color))
	return mainStyle.Render(main.compositor.Render())
}


func (main *MainPane) SetSize() {
	grid := utils.ResolveGrid(main.Dimensions, [][]types.GridCell{
		{{Config: main.requestPane.PaneConfig}},
		{{Config: main.responsePane.PaneConfig}},
	})
	main.requestPane.Dimensions = grid[0][0]
	main.responsePane.Dimensions = grid[1][0]

	main.requestPane.SetSize()
	main.responsePane.SetSize()
}


//event loop listener for: when a new file is to be loaded -> load onto -> url, type, headers, ... 
// when a message sent to stateLoadChan, Update() will be called and this msg is forwarded there
func (main *MainPane) __waitForStateLoad() tea.Cmd {
	return func() tea.Msg {
		return <-main.stateLoadChan 	
	}
}
func (main *MainPane) __stateLoadPerform(state types.HttyState){
	utils.Debugf("State changing requested!")
	main.requestPane.method.SetValue(state.Method)
	main.requestPane.url.SetValue(state.URL)
	main.requestPane.headers.SetValue(utils.HeaderKVEncoder(state.ReqHeaders))
	main.requestPane.body.SetValue(state.ReqBody)
	main.responsePane.SetResponse(
		state.Response.Output, state.Response.Raw, state.Response.Headers, state.Response.Status,
	)
}
func (main *MainPane) __saveStateDialog(){
	utils.Debugf("saving now!")
	main.currStateBuf = types.HttyState{
		Method: main.requestPane.method.GetValue(),
		URL: main.requestPane.url.GetValue(),
		ReqBody: main.requestPane.body.GetValue(),
		ReqHeaders: utils.HeaderKVparser(main.requestPane.headers.GetValue()),
		Response: main.currStateBuf.Response,
	}	
	savePath , err := utils.SaveFileDialog("response.hstate") 
	err = utils.WriteObjectIntoFile(savePath, main.currStateBuf)
	if err != nil {
		panic(err)
	}
}
