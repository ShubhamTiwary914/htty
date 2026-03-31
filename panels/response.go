package htty

import (
	"strings"
	"fmt"
	components "htty/panels/components"
	types "htty/types"
	utils "htty/utils"
	"slices"
	"time"

	global "htty/globals"
	"os"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/browser"
)

type ResponsePane struct {
	textpane components.TextPane
	//objects gets filled after API call from main -> response
	bodyRaw string
	headersRaw map[string]string
	status int

	loading bool
	verboseAllow bool
}

func (res *ResponsePane) Init() tea.Cmd {
	utils.Infof("response panel initialization")
	res.textpane = NewResponseTextComponent()
	res.textpane.Init()
	return nil
}

func (res *ResponsePane) Update(msg tea.Msg) tea.Cmd {
	responseActionKeys := global.Config.Panels.Main_res.Keys		
	actionKeys := utils.GetPanelActionKeys(global.Config.Panels.Main_res)
	focused := (global.CurrentPanelID == global.Config.Panels.Main_res.ID)
	switch msg := msg.(type) {
		case tea.KeyMsg:
		if focused{
			switch msg.String(){
				case responseActionKeys["save"]:
					res.saveDialog()
				case responseActionKeys["open"]:
					res.openView()
				case responseActionKeys["copy"]:	
					res.copyClipboard()
				case responseActionKeys["verbose"]:
					res.verboseToggle()
			}		
			if slices.Contains(actionKeys, msg.String()){
				return nil
			}
		}	
	}
	return res.textpane.Update(msg)
}

func (res ResponsePane) View() string {
	return res.textpane.View()
}

func (res *ResponsePane) SetSize(width int, height int) {
	res.textpane.SetSize(width-2, height)
}

func (res *ResponsePane) SetResponse(formattedResp string, rawBody string, headers map[string]string, status int) {
	res.bodyRaw = rawBody
	res.headersRaw = headers
	res.status = status
	res.textpane.Input.SetValue(formattedResp)
}

func NewResponseTextComponent() components.TextPane {
	var responseTextComponent = components.TextPane{
		CharLimit:   2147483647,
		PanelTitle: global.Config.Panels.Main_res.Title,
		PanelID:     global.PANEL_RES,
		Placeholder: "response will appear here... (API call with ctrl+enter)",
		Showline:    false,
		Border:      types.BorderConfig{Bottom: true, Top: true, Left: true, Right: true},
		StatusOptions: utils.GetPanel_KeyOptions(global.Config.Panels.Main_res),
	}
	return responseTextComponent
}

func (res *ResponsePane) saveDialog() {
	filename, err := utils.SaveFileDialog("response.txt") 
	//may happen if dialog closed force(cancelled), its fine, but is still logged
    if err != nil {
       utils.Errorf("%v", err) 
    }	
	os.WriteFile(filename, []byte(res.textpane.Input.Value()), 0644)
}

//saves body content to temp dir's random file and opens that file in browser/editor/etc...
func (res *ResponsePane) openView() {
	fileType := utils.GetResponseContent_FileExtension(res.headersRaw["Content-Type"])
	randomFileNameUUID := utils.GenerateRandomUUID(16)
	filePath := fmt.Sprintf("%s/%s.%s", global.TEMP_DIR, randomFileNameUUID, fileType)
	content := res.textpane.Input.Value() 
	err := utils.WriteFileContents(filePath, content)
	if err != nil {
		panic(err)
	}
    browser.OpenFile(filePath)
}

func (res *ResponsePane) copyClipboard() {
	clipboard.WriteAll(res.textpane.Input.Value())	
}

func (res *ResponsePane) verboseToggle(){
	res.verboseAllow = !res.verboseAllow
	var newOutput string = utils.ResponseParser_main([]byte(res.bodyRaw), res.headersRaw, 
								res.status, res.verboseAllow) 
	res.SetResponse(newOutput, string(res.bodyRaw), res.headersRaw, res.status)
}

func (res *ResponsePane) HandleLoadTimer(){
	loadTimer := time.NewTicker(1 * time.Second) 
	var dots strings.Builder
	go func () {
		for range loadTimer.C {
			if !res.loading {
				loadTimer.Stop()
				return;
			} 
			utils.Debugf("Timer ticked!!")
			dots .WriteString(".")
			res.SetResponse("Loading" + dots.String(), "", nil, 0)
		}
	}()
}
