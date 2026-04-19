package panels 

import (
	global "htty/globals"
	components "htty/panels/components"
	types "htty/types"
	utils "htty/utils"

	"charm.land/lipgloss/v2"
	tea "github.com/charmbracelet/bubbletea"
)

type RequestPane struct {
	focusIndex    int
	initialized   bool
	method        components.TextOptions
	url           components.TextOptions
	headers       components.TextPane
	body          components.TextPane
	compositor    *lipgloss.Compositor

	Dimensions types.PaneGeometry 
	PaneConfig types.HttyPanel 
}

func (rq *RequestPane) Init() tea.Cmd {
	utils.Infof("request panel initialization")
	rq.PaneConfig = global.Config.Panels.Main_req
	rq.InitSubPanels()	
	return nil
}

func (rq *RequestPane) Update(msg tea.Msg) tea.Cmd {
	return utils.UpdatePanels(msg, &rq.method, &rq.url, &rq.headers, &rq.body)
}


func (rq *RequestPane) View() string {
	rq.buildCompositor()
	reqStyle := lipgloss.NewStyle().Background(lipgloss.Color(global.Config.Common.Background_color))
	return reqStyle.Render(rq.compositor.Render())
}

func (rq *RequestPane) SetSize() {
	rq.method.Dimensions = utils.GetPaneGeometry(rq.method.PaneCfg, rq.Dimensions)
	rq.url.Dimensions = utils.GetPaneGeometry(rq.url.PaneCfg, rq.Dimensions)
	rq.headers.Dimensions = utils.GetPaneGeometry(rq.headers.PaneCfg, rq.Dimensions)
	rq.body.Dimensions = utils.GetPaneGeometry(rq.body.PaneCfg, rq.Dimensions)

	rq.method.SetSize()
	rq.url.SetSize()
	rq.headers.SetSize()
	rq.body.SetSize()
}

//config for all the subpanels for Request
func (rq *RequestPane) InitSubPanels() {
	rq.method = components.TextOptions{
		PanelTitle: rq.method.PaneCfg.Title,
		PanelID:     global.PANEL_REQ_METHOD_ID,
		Placeholder: "Method",
		OptionsFilePath: global.CachePrefix + "/method.txt",
		PaneCfg: global.Config.Panels.Main_req_method,
	}
	rq.url = components.TextOptions{
		PanelTitle: rq.url.PaneCfg.Title,
		PanelID:     global.PANEL_REQ_URL_ID,
		Placeholder: "http://example/com",
		OptionsFilePath: global.CachePrefix + "/url.txt",
		AllowSaveInput: true,
		PaneCfg: global.Config.Panels.Main_req_url,
	}
	rq.headers = components.TextPane{
		PanelTitle: rq.headers.PaneCfg.Title,
		PanelID:     global.PANEL_REQ_HEADERS,
		Placeholder: "Header-Key:   Header-Value\nHeader-Key-2: Header-Value-2\n...",
		PaneCfg: global.Config.Panels.Main_req_headers,
	}
	rq.body = components.TextPane{
		PanelTitle: rq.body.PaneCfg.Title,
		PanelID:     global.PANEL_REQ_BODY,
		Placeholder: "request body content",
		PaneCfg: global.Config.Panels.Main_req_body,
	}
	rq.method.Init()
	rq.url.Init()
	rq.headers.Init()
	rq.body.Init()
}


func (rq *RequestPane) buildCompositor() {
	//request panels offset
	methodWidth := utils.GetPercent(global.Config.Panels.Main_req_method.Width, rq.Dimensions.Width)
	methodHeight := utils.GetPercent(global.Config.Panels.Main_req_method.Height, rq.Dimensions.Height)
	headersWidth := utils.GetPercent(global.Config.Panels.Main_req_headers.Width, rq.Dimensions.Width)

	//collect all layers
	methodLayer := utils.CreateNewLayer(&rq.method, global.Config.Panels.Main_req_method)
	_, methodOptionsLayer := rq.method.ViewWithOptions(true)
	urlLayer := utils.CreateNewLayer(&rq.url, global.Config.Panels.Main_req_url, methodWidth)
	_, urlOptionsLayer := rq.url.ViewWithOptions(true); 

	headersLayer := utils.CreateNewLayer(&rq.headers, global.Config.Panels.Main_req_headers, 0, methodHeight)
	bodyLayer := utils.CreateNewLayer(&rq.body, global.Config.Panels.Main_req_body, headersWidth, methodHeight)	
	layers := []*lipgloss.Layer{methodLayer, urlLayer, headersLayer, bodyLayer}

	// add options overlay if focused
	if methodOptionsLayer != nil {
		methodOptionsLayer.X(3).Y(methodHeight+3)
		layers = append(layers, methodOptionsLayer)
	}
	if urlOptionsLayer != nil {
		urlOptionsLayer.X(global.Config.Panels.Main_req_url.Margin[0]+11).Y(methodHeight+ 3)
		layers = append(layers, urlOptionsLayer)
	}	
	rq.compositor = lipgloss.NewCompositor(layers...)
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

