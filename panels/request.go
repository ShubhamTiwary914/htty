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
	rq.__initSubPanels()	
	return nil
}

func (rq *RequestPane) Update(msg tea.Msg) tea.Cmd {
	return utils.UpdatePanels(msg, &rq.method, &rq.url, &rq.headers, &rq.body)
}


func (rq *RequestPane) View() string {
	rq.__buildCompositor()
	reqStyle := lipgloss.NewStyle().Background(lipgloss.Color(global.Config.Common.Background_color))
	return reqStyle.Render(rq.compositor.Render())
}

func (rq *RequestPane) SetSize() {
	grid := utils.ResolveGrid(rq.Dimensions, [][]types.GridCell{
		{{Config: rq.method.PaneCfg}, {Config: rq.url.PaneCfg}},
		{{Config: rq.headers.PaneCfg}, {Config: rq.body.PaneCfg}},
	})

	rq.method.Dimensions  = grid[0][0]
	rq.url.Dimensions     = grid[0][1]
	rq.headers.Dimensions = grid[1][0]
	rq.body.Dimensions    = grid[1][1]

	rq.method.SetSize()
	rq.url.SetSize()
	rq.headers.SetSize()
	rq.body.SetSize()
}

//config for all the subpanels for Request
func (rq *RequestPane) __initSubPanels() {
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


func (rq *RequestPane) __buildCompositor() {
	//collect all layers
	methodLayer := utils.CreateLayerFromDims(&rq.method, rq.method.Dimensions, 1)
	urlLayer := utils.CreateLayerFromDims(&rq.url, rq.url.Dimensions, 1)
	headersLayer := utils.CreateLayerFromDims(&rq.headers, rq.headers.Dimensions, 1)
	bodyLayer := utils.CreateLayerFromDims(&rq.body, rq.body.Dimensions, 1)
	layers := []*lipgloss.Layer{methodLayer, urlLayer, headersLayer, bodyLayer}

	_, methodOptionsLayer := rq.method.ViewWithOptions(true)
	_, urlOptionsLayer := rq.url.ViewWithOptions(true); 
	// add options overlay if focused
	if methodOptionsLayer != nil {
		methodOptionsLayer.X(rq.method.Dimensions.X).Y(rq.method.Dimensions.Y + rq.method.Dimensions.Height+1).Z(2)
		layers = append(layers, methodOptionsLayer)
	}
	if urlOptionsLayer != nil {
		urlOptionsLayer.X(rq.url.Dimensions.X).Y(rq.url.Dimensions.Y + rq.url.Dimensions.Height+1).Z(2)
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

