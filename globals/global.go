package htty

import (
	types "htty/types"
	"os"
	"log"
)

var CurrentPanelID = types.PANEL_FOCUS_IDS[types.PANEL_REQ_ID]
var Logger = log.New(os.Stdout, "", 0)

