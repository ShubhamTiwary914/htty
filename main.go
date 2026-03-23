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
		utils.Errorf("%v", err)
		panic(err)			
	}
	//overwrite = true: flushes logfile each run
	global.LOGLEVEL = os.Getenv("LOGLEVEL")
	utils.RedirectLogs_toFile(global.Config.Log.File, true)
	utils.Infof("htty application initializing")
	//envs set
	global.CachePrefix = os.Getenv("CACHE_PREFIX")
	global.PANEL_FOCUS_IDS, err = utils.GetPanelIDsMap(global.Config)
	global.CurrentPanelID = global.PANEL_FOCUS_IDS[global.PANEL_REQ_METHOD_ID]
	if err != nil {
		panic(err)			
	}
	global.TEMP_DIR = os.Getenv("TMP_DIR")
	utils.Infof("Env Keys: loglevel(%s), cache_prefix(%s), temp_dir(%s)", global.LOGLEVEL, 
		global.CachePrefix, global.TEMP_DIR, 
	)
	//create tea program object
	app := App{}
	return app
}
