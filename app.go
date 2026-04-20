package main

import (
	global "htty/globals"
	panels "htty/panels"
	components "htty/panels/components"
	"htty/types"
	utils "htty/utils"

	"charm.land/lipgloss/v2"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	sidePane   			panels.SidePane
	mainPane   			panels.MainPane
	statusLinePane 		panels.StatusLinePane
	alertPane     		components.AlertPane
	compositor 			*lipgloss.Compositor

	alertBoxChan chan any 
	Dimensions          types.PaneGeometry
}

func (app *App) Init() tea.Cmd {
	utils.Infof("app panel initialization called")
	events := app.__eventsInit()
	mainTea := app.mainPane.Init()
	sideTea := app.sidePane.Init()
	statusTea := app.statusLinePane.Init()
	return tea.Batch(events, mainTea, sideTea, statusTea)
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	
	//auto resize to window dimensions
	case tea.WindowSizeMsg:
		global.AppWidth = msg.Width
		global.AppHeight = msg.Height
		app.SetSize()
		app.View() 
	case tea.KeyMsg:
		if msg.String() == global.Config.Key.Quit {
			utils.Infof("exiting the htty program...")
			return &app, tea.Quit
		}
		if msg.String() == global.Config.Key.Nextpanel {
			utils.PanelFocusNext(&global.CurrentPanelID)
		}
		//jump sequence
		allowJump, targetPanel := utils.EventIs_TypeJumpPanel(msg.String()); if allowJump {
			if targetPanel < global.FOCUSABLE_PANELS{
				utils.PanelFocusJump(&global.CurrentPanelID, targetPanel)
				return &app, nil
			}
		}
	case components.AlertPane:
		return app.__alertboxHandler(msg)
	}
	//allows passing tea object for handling events to children panes
	return &app, utils.UpdatePanels(msg, &app.sidePane, &app.mainPane, &app.alertPane)
}


func (app App) View() string {
	sideLayer := utils.CreateLayerFromDims(&app.sidePane, app.sidePane.Dimensions, 1)
	mainLayer := utils.CreateLayerFromDims(&app.mainPane, app.mainPane.Dimensions, 1)
	statusLineLayer := utils.CreateLayerFromDims(&app.statusLinePane, app.statusLinePane.Dimensions, 1)
	
	layers := []*lipgloss.Layer{sideLayer, mainLayer, statusLineLayer}
	if alertLayer := app.alertPane.ViewAsLayer(); alertLayer != nil {
		layers = append(layers, alertLayer)
	}

	app.compositor = lipgloss.NewCompositor(layers...)
	appStyle := lipgloss.NewStyle().Background(lipgloss.Color(global.Config.Common.Background_color))
	return appStyle.Render(app.compositor.Render())
}


func (app *App) SetSize() {
	app.Dimensions = types.PaneGeometry{Width: global.AppWidth, Height: global.AppHeight}
	grid := utils.ResolveGrid(app.Dimensions, [][]types.GridCell{
		{{Config: app.sidePane.PaneConfig}, {Config: app.mainPane.PaneConfig}},
		{{Config: app.statusLinePane.PaneConfig}},
	})
	app.sidePane.Dimensions       = grid[0][0]
	app.mainPane.Dimensions       = grid[0][1]
	app.statusLinePane.Dimensions = grid[1][0]

	app.mainPane.SetSize()
	app.sidePane.SetSize()
}


//INFO: EVENTS HANDLER section (app is top level panel and can handle IO, so lots may be appear here)

//make channels and register as subscribers
func (app *App) __eventsInit() (tea.Cmd){
	app.alertBoxChan = make(chan any, 1)
	global.StateBus.Subscribe(global.EVENT_ALERTPANE, app.alertBoxChan)
	utils.Infof("app subscribed to event %s", global.EVENT_ALERTPANE)

	alertBoxWait := func() tea.Msg {
		return <- app.alertBoxChan 		
	}
	return tea.Batch(alertBoxWait)
}

func (app *App) __alertboxHandler(alertPaneCfg components.AlertPane) (tea.Model, tea.Cmd) {
	utils.Debugf("app has received the event %s!", global.EVENT_ALERTPANE)

	app.alertPane.Dimensions = alertPaneCfg.Dimensions
	app.alertPane.TTL        = alertPaneCfg.TTL
	app.alertPane.EndKey     = alertPaneCfg.EndKey
	showCmd := app.alertPane.Show(alertPaneCfg.Message, alertPaneCfg.Level)

	// re-arm: keep listening for the next alert event
	rearmCmd := func() tea.Msg {
		return <-app.alertBoxChan
	}
	return app, tea.Batch(showCmd, rearmCmd)
}
