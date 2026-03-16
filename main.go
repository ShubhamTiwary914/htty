package main

import (
	global "htty/globals"
	utils "htty/utils"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main(){
	app := __init()
	proc := tea.NewProgram(&app, tea.WithAltScreen())
	if _, err := proc.Run(); err != nil {
        utils.Errorf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}

func __init() App {
	//load config
	var err error
	global.Config , err = utils.GetConfig() 
	if err != nil {
		panic(err)			
	}
	global.PANEL_FOCUS_IDS, err = utils.GetPanelIDsMap(global.Config)
	global.CurrentPanelID = global.PANEL_FOCUS_IDS[global.PANEL_REQ_METHOD_ID]
	if err != nil {
		panic(err)			
	}
	//overwrite = true: flushes logfile each run
	utils.RedirectLogs_toFile(global.Config.Log.File, true)
	utils.Infof("htty application initializing")
	//create tea program object
	app := App{}
	return app
}
