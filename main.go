package main

import (
	"fmt"
	global "htty/globals"
	utils "htty/utils"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main(){
	__init()

	app := App{}
	proc := tea.NewProgram(&app, tea.WithAltScreen())
	if _, err := proc.Run(); err != nil {
        utils.Errorf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}

func __init() {
	//load config
	var err error
	homeDir, err := os.UserHomeDir()
	if err != nil {
		utils.Errorf("%v", err)
		panic(err)
	}
	global.CONFIG_PATH = utils.GetEnv("CONFIG_FILE", fmt.Sprintf(
		"/home/%s/.config/htty/config.json", homeDir,
	))
	global.Config , err = utils.GetConfig() 
	if err != nil {
		utils.Errorf("%v", err)
		panic(err)			
	}

	//overwrite = true: flushes logfile each run
	global.LOGLEVEL = utils.GetEnv("LOGLEVEL", "info")
	global.LOGFILE = utils.GetEnv("LOGFILE", "/var/log/htty/htty.log")
	utils.RedirectLogs_toFile(global.LOGFILE, true)
	utils.Infof("htty application initializing")

	//envs set
	global.CachePrefix = utils.GetEnv("CACHE_PREFIX", fmt.Sprintf(
		"/home/%s/.cache/htty", homeDir,
	))
	global.PANEL_FOCUS_IDS, err = utils.GetPanelIDsMap(global.Config)
	global.CurrentPanelID = global.PANEL_FOCUS_IDS[global.PANEL_REQ_METHOD_ID]
	if err != nil {
		panic(err)			
	}
	global.TEMP_DIR = utils.GetEnv("TMP_DIR", "/tmp")
	utils.Infof("Env Config: loglevel(%s), cache_prefix(%s), temp_dir(%s)", 
		global.LOGLEVEL, global.CachePrefix, global.TEMP_DIR, 
	)
}
