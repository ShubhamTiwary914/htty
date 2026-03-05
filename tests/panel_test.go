package htty_test

import (
	utils "htty/utils"	
	types "htty/types"
	"testing"
)

func TestPanelFocusControl(tt *testing.T){
	var currentPanel int = 0	
	//forwards
	utils.PanelFocusNext(&currentPanel)
	if currentPanel != 1 {
		tt.Errorf("panelID: %d, expected: %d", currentPanel, 1)
	}
	//back
	utils.PanelFocusPrev(&currentPanel)
	if currentPanel != 0 {
		tt.Errorf("panelID: %d, expected: %d", currentPanel, 0)
	}
	//jumps
	utils.PanelFocusJump(&currentPanel, 2)
	if currentPanel != 2 {
		tt.Errorf("panelID: %d, expected: %d", currentPanel, 2)
	}

	utils.PanelFocusJump(&currentPanel, types.PANEL_FOCUS_IDS["SIDE"])
	if currentPanel != 0 {
		tt.Errorf("panelID: %d, expected: %d", currentPanel, 0)
	}
}

