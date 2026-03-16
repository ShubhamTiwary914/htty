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

const (
	PANEL_SIDE_ID string = "side"
	PANEL_REQ_ID string = "main_req"
	PANEL_REQ_METHOD_ID string = "main_req_method"
	PANEL_REQ_URL_ID string = "main_req_url"
	PANEL_REQ_HEADERS string = "main_req_headers"
	PANEL_REQ_BODY string = "main_req_body"
	PANEL_RES string = "main_res"
)

//enum for "focused" panel where cmd actions can act on currently
var PANEL_FOCUS_IDS = map[string]int{}

var CurrentPanelID = PANEL_FOCUS_IDS[PANEL_REQ_ID]
var Logger = log.New(os.Stdout, "", 0)
var Config types.HttyConfig

var AppWidth int
var AppHeight int

