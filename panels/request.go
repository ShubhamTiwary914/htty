package htty

import (
	components "htty/panels/components"
	types "htty/types"
	utils "htty/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)		

type RequestPane struct {
	width, height int
	focusIndex int
	initialized bool
	method components.TextPane 
	url components.TextPane
	headers components.TextPane 
	body components.TextPane 
}


func (rq *RequestPane) Init() tea.Cmd {
	utils.Infof("request panel initialization")
	rq.method, rq.url, rq.headers, rq.body = RequestSubPanels()
	rq.method.Init()
	rq.url.Init()
	rq.headers.Init()
	rq.body.Init()
	return nil
}

func (rq *RequestPane) Update(msg tea.Msg) (tea.Cmd){	
	return utils.UpdatePanels(msg, &rq.method, &rq.url, &rq.headers, &rq.body)
}

func (rq RequestPane) View() string {
	style := utils.SetFullBorder(rq.width-2, rq.height, lipgloss.Color(utils.GetPanelFocusColor(types.PANEL_REQ_ID))) 
	firstRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		rq.method.View(),
		rq.url.View(),
	)
	secondRaw := lipgloss.JoinHorizontal(
		lipgloss.Top,
		rq.headers.View(),
		rq.body.View(),
	)
	requestSubpanels := lipgloss.JoinVertical(
		lipgloss.Left,
		firstRow,
		secondRaw,
	)
	return style.Render(requestSubpanels)
}

func (rq *RequestPane) SetSize(w int, h int){
	//TODO: instead of raw dogging sizes, its better to load based off config, with relative sizing in %s 
	rq.width = w
	rq.height = h
	rq.method.SetSize(w/11, h/12)
	rq.url.SetSize(int(float64(w)/1.25), h/12)
	rq.headers.SetSize(w/3, int(float64(h)/1.5))
	rq.body.SetSize(int(float64(w)/1.8), int(float64(h)/1.5))
}

//config for all the subpanels for Request
func RequestSubPanels() (components.TextPane, components.TextPane, components.TextPane, components.TextPane) {
	//TODO: these config are better off handled by the config manager instead of this
	var methodTypeComponent = components.TextPane{
		CharLimit: 10, PanelID: types.PANEL_REQ_METHOD_ID, 
		Placeholder: "Method", Showline: false,
		Border: types.BorderConfig{Bottom: true},
		Margin: types.MarginConfig{Left:3, Top: 1},
	}
	var urlPathComponent = components.TextPane{
		CharLimit: 1024, PanelID: types.PANEL_REQ_URL_ID,
		Placeholder: "http://example/com", Showline: false,
		Border: types.BorderConfig{Bottom: true, Top: true, Left: true, Right: true},
		Margin: types.MarginConfig{Left: 5},
	}
	var headersComponent = components.TextPane{
		CharLimit: 1024,
		PanelID: types.PANEL_REQ_HEADERS, Placeholder: "Header-Key:   Header-Value\nHeader-Key-2: Header-Value-2\n...",
		Showline: true,
		Border: types.BorderConfig{Bottom: true, Top: true, Left: true, Right: true},
		Margin: types.MarginConfig{Left: 3, Top: 1},
	}
	var bodyComponent = components.TextPane{
		CharLimit: 2048,
		PanelID: types.PANEL_REQ_BODY, Placeholder: "request body content",
		Showline: true,
		Border: types.BorderConfig{Bottom: true, Top: true, Left: true, Right: true},
		Margin: types.MarginConfig{Left: 4, Top: 1},
	}
	return methodTypeComponent, urlPathComponent, headersComponent, bodyComponent
}


//take current snapshot of inputs: method,url,headers,body & compose a HttpType object
func (rq *RequestPane) ExportPayload() (types.HttpType) {
	method := rq.method.Input.Value()
	path := rq.url.Input.Value()	
	body:= rq.body.Input.Value()

	rawHeaders := rq.headers.Input.Value()
	headers := utils.HeaderKVparser(rawHeaders)	

	return types.HttpType{
		Method: method,
		Path: path,
		Headers: headers,
		Body: body,
	}
}
