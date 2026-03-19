package htty

import (
	global "htty/globals"
	components "htty/panels/components"
	types "htty/types"
	utils "htty/utils"

	"charm.land/lipgloss/v2"
	tea "github.com/charmbracelet/bubbletea"
)

type RequestPane struct {
	width, height int
	focusIndex    int
	initialized   bool
	method        components.TextOptions
	url           components.TextPane
	headers       components.TextPane
	body          components.TextPane
	compositor    *lipgloss.Compositor
}

func (rq *RequestPane) Init() tea.Cmd {
	utils.Infof("request panel initialization")
	rq.method, rq.url, rq.headers, rq.body = RequestSubPanels()
	rq.method.OptionsStore, _ = utils.ReadTextLines_intoList(global.CachePrefix + "/method.txt")
	rq.method.Init()

	rq.url.Init()
	rq.headers.Init()
	rq.body.Init()
	return nil
}

func (rq *RequestPane) Update(msg tea.Msg) tea.Cmd {
	return utils.UpdatePanels(msg, &rq.method, &rq.url, &rq.headers, &rq.body)
}

func (rq *RequestPane) buildCompositor() {
	//request panels offset
	methodWidth := utils.GetPercent(global.Config.Panels.Main_req_method.Width, rq.width)
	methodHeight := utils.GetPercent(global.Config.Panels.Main_req_method.Height, rq.height)
	headersWidth := utils.GetPercent(global.Config.Panels.Main_req_headers.Width, rq.width)

	//collect all layers
	methodLayer := utils.CreateNewLayer(&rq.method, global.Config.Panels.Main_req_method)
	_, methodOptionsLayer := rq.method.ViewWithOptions(true)
	urlLayer := utils.CreateNewLayer(&rq.url, global.Config.Panels.Main_req_url, methodWidth)
	headersLayer := utils.CreateNewLayer(&rq.headers, global.Config.Panels.Main_req_headers, 0, methodHeight)
	bodyLayer := utils.CreateNewLayer(&rq.body, global.Config.Panels.Main_req_body, headersWidth, methodHeight)	
	layers := []*lipgloss.Layer{methodLayer, urlLayer, headersLayer, bodyLayer}

	// add options overlay if focused
	if methodOptionsLayer != nil {
		methodOptionsLayer.X(3).Y(methodHeight + 2)
		layers = append(layers, methodOptionsLayer)
	}
	rq.compositor = lipgloss.NewCompositor(layers...)
}

func (rq *RequestPane) View() string {
	rq.buildCompositor()
	return rq.compositor.Render()
}

func (rq *RequestPane) SetSize(width int, height int) {
	rq.width = width
	rq.height = height
	rq.method.SetSize(
		utils.GetPercent(global.Config.Panels.Main_req_method.Width, rq.width),
		utils.GetPercent(global.Config.Panels.Main_req_method.Height, rq.height),
	)
	rq.url.SetSize(
		utils.GetPercent(global.Config.Panels.Main_req_url.Width, rq.width),
		utils.GetPercent(global.Config.Panels.Main_req_url.Height, rq.height),
	)
	rq.headers.SetSize(
		utils.GetPercent(global.Config.Panels.Main_req_headers.Width, rq.width),
		utils.GetPercent(global.Config.Panels.Main_req_headers.Height, rq.height),
	)
	rq.body.SetSize(
		utils.GetPercent(global.Config.Panels.Main_req_body.Width, rq.width),
		utils.GetPercent(global.Config.Panels.Main_req_body.Height, rq.height),
	)
}

//config for all the subpanels for Request
func RequestSubPanels() (components.TextOptions, components.TextPane, components.TextPane, components.TextPane) {
	var methodTypeComponent = components.TextOptions{
		CharLimit:   10,
		PanelID:     global.PANEL_REQ_METHOD_ID,
		Placeholder: "Method",
		Showline:    false,
		Border:      types.BorderConfig{Bottom: true},
	}
	var urlPathComponent = components.TextPane{
		CharLimit:   1024,
		PanelID:     global.PANEL_REQ_URL_ID,
		Placeholder: "http://example/com",
		Showline:    false,
		Border:      types.BorderConfig{Bottom: true},
	}
	var headersComponent = components.TextPane{
		CharLimit:   1024,
		PanelID:     global.PANEL_REQ_HEADERS,
		Placeholder: "Header-Key:   Header-Value\nHeader-Key-2: Header-Value-2\n...",
		Showline:    true,
		Border:      types.BorderConfig{Bottom: true, Top: true, Left: true, Right: true},
		Margin:      types.MarginConfig{Bottom: 2},
	}
	var bodyComponent = components.TextPane{
		CharLimit:   2048,
		PanelID:     global.PANEL_REQ_BODY,
		Placeholder: "request body content",
		Showline:    true,
		Border:      types.BorderConfig{Bottom: true, Top: true, Left: true, Right: true},
		Margin:      types.MarginConfig{Left: 0, Bottom: 2, Right: 10},
	}
	return methodTypeComponent, urlPathComponent, headersComponent, bodyComponent
}

//take current snapshot of inputs: method,url,headers,body & compose a HttpType object
func (rq *RequestPane) ExportPayload() types.HttpType {
	method := rq.method.Input.Value()
	path := rq.url.Input.Value()
	body := rq.body.Input.Value()

	rawHeaders := rq.headers.Input.Value()
	headers := utils.HeaderKVparser(rawHeaders)

	return types.HttpType{
		Method:  method,
		Path:    path,
		Headers: headers,
		Body:    body,
	}
}
