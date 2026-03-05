package main

import (
	panels "htty/panels"
	utils "htty/utils"
	global "htty/globals"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	width int 
	height int
	sidePane panels.SidePane 
	mainPane panels.MainPane
}

func (app *App) Init() tea.Cmd {
	utils.Infof("app panel initialization called")
	app.mainPane.Init()
	return nil
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	//auto resize to window dimensions
	case tea.WindowSizeMsg:
		app.SetSize(msg.Width, msg.Height)	
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			utils.Debugf("exting the htty program...")
			return &app, tea.Quit
		}
		if msg.String() == "tab" {
			utils.PanelFocusNext(&global.CurrentPanelID)				
		}
	}	
	//INFO: allows passing tea object for handling events to children panes
	return &app, utils.UpdatePanels(msg, &app.sidePane, &app.mainPane) 
}

func (app App) View() string {
	if app.width == 0 || app.height == 0 {
		return ""
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		app.sidePane.View(),
		app.mainPane.View(),
	)
}

func (app *App) SetSize(width int, height int) {
	app.width = width 
	app.height = height 
	//set children bounds 
	sideWidth := 30
	mainWidth := app.width - sideWidth
	app.sidePane.SetSize(sideWidth, app.height, 2)
	app.mainPane.SetSize(mainWidth, app.height, 2)
}
