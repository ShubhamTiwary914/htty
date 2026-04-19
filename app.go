package main

import (
	global "htty/globals"
	panels "htty/panels"
	"htty/types"
	utils "htty/utils"

	"charm.land/lipgloss/v2"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	sidePane   			panels.SidePane
	mainPane   			panels.MainPane
	statusLinePane 		panels.StatusLinePane
	compositor 			*lipgloss.Compositor

	Dimensions          types.PaneGeometry
}

func (app *App) Init() tea.Cmd {
	utils.Infof("app panel initialization called")
	mainTea := app.mainPane.Init()
	sideTea := app.sidePane.Init()
	statusTea := app.statusLinePane.Init()
	return tea.Batch(mainTea, sideTea, statusTea)
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
	}
	//allows passing tea object for handling events to children panes
	return &app, utils.UpdatePanels(msg, &app.sidePane, &app.mainPane)
}


func (app App) View() string {
	sideLayer := utils.CreateLayerFromDims(&app.sidePane, app.sidePane.Dimensions, 1)
	mainLayer := utils.CreateLayerFromDims(&app.mainPane, app.mainPane.Dimensions, 1)
	statusLineLayer := utils.CreateLayerFromDims(&app.statusLinePane, app.statusLinePane.Dimensions, 1)

	app.compositor = lipgloss.NewCompositor(sideLayer, mainLayer, statusLineLayer)
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
