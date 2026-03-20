package htty

import (
	components "htty/panels/components"
	types "htty/types"
	utils "htty/utils"

	global "htty/globals"
	tea "github.com/charmbracelet/bubbletea"
)

type ResponsePane struct {
	textpane components.TextPane
}

func (res *ResponsePane) Init() tea.Cmd {
	utils.Infof("response panel initialization")
	res.textpane = NewResponseTextComponent()
	res.textpane.Init()
	return nil
}

func (res *ResponsePane) Update(msg tea.Msg) tea.Cmd {
	return res.textpane.Update(msg)
}

func (res ResponsePane) View() string {
	return res.textpane.View()
}

func (res *ResponsePane) SetSize(width int, height int) {
	res.textpane.SetSize(width-2, height)
}

func (res *ResponsePane) SetResponse(body string) {
	res.textpane.Input.SetValue(body)
}

func NewResponseTextComponent() components.TextPane {
	var responseTextComponent = components.TextPane{
		CharLimit:   2147483647,
		PanelID:     global.PANEL_RES,
		Placeholder: "response will appear here... (API call with ctrl+enter)",
		Showline:    false,
		Border:      types.BorderConfig{Bottom: true, Top: true, Left: true, Right: true},
		Margin:      types.MarginConfig{Left: 1, Top: 1},
	}
	return responseTextComponent
}
