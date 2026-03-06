/*
	Variables that are to be persisted across files, that are very much needed in almost every single function
	NOTE: this is to be used as least as possible, for any constant, save under types/ (makes test/mocks much easier to replicate)
*/

package htty

import (
	types "htty/types"
	"log"
	"os"
)

var CurrentPanelID = types.PANEL_FOCUS_IDS[types.PANEL_REQ_ID]
var Logger = log.New(os.Stdout, "", 0)
var Config types.HttyConfig

var AppWidth int
var AppHeight int
