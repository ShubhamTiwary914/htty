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

//enum for "focused" panel where cmd actions can act on currently
var PANEL_FOCUS_IDS = map[string]int{}

var CurrentPanelID = PANEL_FOCUS_IDS[PANEL_REQ_ID]
var Logger = log.New(os.Stdout, "", 0)
var Config types.HttyConfig

var CachePrefix string
var AppWidth int
var AppHeight int

var HTTP_METHOD = map[string]int{
	"GET": 1,
	"HEAD": 2,
	"POST": 3,
	"PUT": 4,
	"PATCH": 5,
	"DELETE": 6,
	"CONNECT": 7,
	"OPTIONS": 8,
	"TRACE": 9,
}

const (
	PANEL_SIDE_ID string = "side"
	PANEL_REQ_ID string = "main_req"
	PANEL_REQ_METHOD_ID string = "main_req_method"
	PANEL_REQ_URL_ID string = "main_req_url"
	PANEL_REQ_HEADERS string = "main_req_headers"
	PANEL_REQ_BODY string = "main_req_body"
	PANEL_RES string = "main_res"
)

const (
	REQUEST_TIMEOUT int = 10
	FOCUSABLE_PANELS int  = 6
)

const (
	LOG_INFO  = "info"
	LOG_WARN  = "warn"
	LOG_ERROR = "error"
	LOG_DEBUG = "debug"
	LOG_ALL   = "all"
	LOG_ENVNAME = "LOGLEVEL"
)

const (
	BORDER_UP    = "BORDER_UP"
	BORDER_DOWN  = "BORDER_DOWN"
	BORDER_LEFT  = "BORDER_LEFT"
	BORDER_RIGHT = "BORDER_RIGHT"
)
