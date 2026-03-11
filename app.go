package main

import (
	global "htty/globals"
	panels "htty/panels"
	utils "htty/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	sidePane panels.SidePane 
	mainPane panels.MainPane
}

func (app *App) Init() tea.Cmd {
	utils.Infof("app panel initialization called")
	app.mainPane.Init()
	app.sidePane.Init()
	return nil
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	//auto resize to window dimensions
	case tea.WindowSizeMsg:
		global.AppWidth = msg.Width
		global.AppHeight = msg.Height
		app.SetSize()	
	case tea.KeyMsg:
		if msg.String() == global.Config.Key.Quit {
			utils.Debugf("exting the htty program...")
			return &app, tea.Quit
		}
		if msg.String() == global.Config.Key.Nextpanel {
			utils.PanelFocusNext(&global.CurrentPanelID)				
		}
	}	
	//INFO: allows passing tea object for handling events to children panes
	return &app, utils.UpdatePanels(msg, &app.sidePane, &app.mainPane) 
}

func (app App) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		app.sidePane.View(),
		app.mainPane.View(),
	)
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
