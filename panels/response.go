/*
	wraps around text component but seperated since it probably needs to be extended off due to its importance
	present as child of main panel
*/

package htty

import (
	components "htty/panels/components"
	types "htty/types"
	utils "htty/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type ResponsePane struct {
	width, height int
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

func (res *ResponsePane) SetSize(w, h int) {
	res.width = w
	res.height = h
	res.textpane.SetSize(int(float64(w)/1.05), h)
}

func (res *ResponsePane) SetResponse(body string) {
	res.textpane.Input.SetValue(body)
}

func NewResponseTextComponent() (components.TextPane) {
	var responseTextComponent = components.TextPane{
		CharLimit: 1024, PanelID: types.PANEL_RES, 
		Placeholder: "response will appear here... (API call with ctrl+enter)", Showline: false,
		Border: types.BorderConfig{Bottom: true, Top: true, Left: true, Right: true},
		Margin: types.MarginConfig{Left:1,Top: 1},
	}
	return responseTextComponent 
}
