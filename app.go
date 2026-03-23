package main

import (
	global "htty/globals"
	panels "htty/panels"
	utils "htty/utils"

	"charm.land/lipgloss/v2"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	sidePane   			panels.SidePane
	mainPane   			panels.MainPane
	statusLinePane 		panels.StatusLinePane
	compositor 			*lipgloss.Compositor
}

func (app *App) Init() tea.Cmd {
	utils.Infof("app panel initialization called")
	app.mainPane.Init()
	app.sidePane.Init()
	app.statusLinePane.Init()
	app.SetSize()
	app.View()
	return nil
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
	//INFO: allows passing tea object for handling events to children panes
	return &app, utils.UpdatePanels(msg, &app.sidePane, &app.mainPane)
}


func (app App) View() string {
	//set the side and main layers 
	sideLayer := utils.CreateNewLayer(&app.sidePane, global.Config.Panels.Side)
	mainLayer := utils.CreateNewLayer(&app.mainPane, global.Config.Panels.Main, 
		utils.GetPercent(global.Config.Panels.Side.Width, global.AppWidth),
	)	
	statusLineLayer := utils.CreateNewLayer(&app.statusLinePane, global.Config.Panels.Statusline,
		0, utils.GetPercent(global.Config.Panels.Main.Height, global.AppHeight),
	)
	app.compositor = lipgloss.NewCompositor(sideLayer, mainLayer, statusLineLayer)
	appStyle := lipgloss.NewStyle().Background(lipgloss.Color(global.Config.Common.Background_color))
	return appStyle.Render(app.compositor.Render())
}

func (app *App) SetSize() {
	app.sidePane.SetSize(
		utils.GetPercent(global.Config.Panels.Side.Width, global.AppWidth),
		utils.GetPercent(global.Config.Panels.Side.Height, global.AppHeight),
	)
	app.mainPane.SetSize(
		utils.GetPercent(global.Config.Panels.Main.Width, global.AppWidth),
		utils.GetPercent(global.Config.Panels.Main.Height, global.AppHeight),
	)
}
