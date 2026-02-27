package htty

import (
	types "htty/types"
)


func PanelFocusNext(focusID *int){
	(*focusID)++
}
func PanelFocusPrev(focusID *int){
	(*focusID)--
}

//jump to panel using the panelID or panel mapping string in types/panels/PANEL_FOCUS
func PanelFocusJump(focusID *int, newfocuskey interface{}){
	switch newfocuskey.(type) {
		case int:
			*focusID = newfocuskey.(int)
		case string:
			*focusID = (types.PANEL_FOCUS[(newfocuskey.(string))])
	}	
}
