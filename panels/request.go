package htty

import (
	types "htty/types"
	utils "htty/utils"
	components "htty/panels/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RequestPane struct {
	width, height int
	focusIndex int
	initialized bool
	textdrop components.TextModel 
}

func (rq *RequestPane) Init() tea.Cmd {
	utils.Infof("request subpanel initialization")
	rq.textdrop = components.TextModel{
		CharLimit: 10, PanelID: types.PANEL_REQ_METHOD_ID, 
		Placeholder: "Method", Showline: false,
		Border: types.BorderConfig{Bottom: true},
	}
	return rq.textdrop.Init()
}

func (rq *RequestPane) Update(msg tea.Msg) (tea.Cmd){	
	return utils.UpdatePanels(msg, &rq.textdrop)
}

func (rq RequestPane) View() string {
	style := utils.SetFullBorder(rq.width-2, rq.height, 
		lipgloss.Color(utils.GetPanelFocusColor(types.PANEL_REQ_ID))) 
	return style.Render(rq.textdrop.View())
}

func (rq *RequestPane) SetSize(w int, h int){
	rq.width = w
	rq.height = h
	rq.textdrop.SetSize(w/11, h/12)
}
