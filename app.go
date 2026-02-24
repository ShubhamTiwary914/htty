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
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	//auto resize to window dimensions
	case tea.WindowSizeMsg:
		app.width = msg.Width
		app.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return app, tea.Quit
		}
	}

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
	side, main := _appstyle(app)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		side,
		main,
	)
}

func _appstyle(app App) (side string, main string){
	borderSize := 2 
	leftWidth :=  app.getWidthSize(30)  - borderSize  
	rightWidth := app.getWidthSize(70) - borderSize 
	height := app.height - borderSize

	frame := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Height(height)
	sideframe := frame.Width(leftWidth).Render(app.sidePane.View())
	mainframe := frame.Width(rightWidth).Render(app.mainPane.View())
	return sideframe, mainframe
}

func (app App) getWidthSize(percent int) int {
	return (app.width * percent/100) 
}
