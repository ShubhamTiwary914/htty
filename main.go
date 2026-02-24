package main

import (
	"fmt"
	panels "htty/panels"
	utils "htty/utils"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

 
func main(){
	if os.Getenv("MODE") == "debug" {
		utils.RedirectLogs_toFile("debug.log", true)
	}
	utils.Debugf("htty application has started...")

	app := App{
		sidePane: panels.SidePane{},
		mainPane: panels.MainPane{},
	}
	proc := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := proc.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}
