package main

import (
	panels "htty/panels"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	width int 
	height int
	sidePane panels.SidePane 
	mainPane panels.MainPane
}

func (app App) Init() tea.Cmd {
	return nil
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	//auto resize to window dimensions
	case tea.WindowSizeMsg:
		app.width = msg.Width
		app.height = msg.Height		
		//set children bounds 
		sideWidth := 30
		mainWidth := app.width - sideWidth
		app.sidePane.SetSize(sideWidth, app.height, 2)
		app.mainPane.SetSize(mainWidth, app.height, 2)

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return app, tea.Quit
		}
	}
	
	var cmds []tea.Cmd
	var cmd tea.Cmd
	app.sidePane, cmd = app.sidePane.Update(msg)
	cmds = append(cmds, cmd)

	app.mainPane, cmd = app.mainPane.Update(msg)
	cmds = append(cmds, cmd)
	return app, tea.Batch(cmds...)
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
