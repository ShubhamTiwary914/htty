package htty

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	PANEL_SIDE_ID string = "SIDE"

	PANEL_REQ_ID string = "MAIN"
	PANEL_REQ_METHOD_ID string = "REQ.METHOD"
	PANEL_REQ_URL_ID string = "REQ.URL"
	PANEL_REQ_HEADERS string = "REQ.HEADERS"
	PANEL_REQ_BODY string = "REQ.BODY"

	PANEL_RES string = "RES"
)
//enum for "focused" panel where cmd actions can act on currently
var PANEL_FOCUS_IDS = map[string]int{
	PANEL_SIDE_ID: 0,
	PANEL_REQ_ID: 1,
	PANEL_REQ_METHOD_ID: 2,
	PANEL_REQ_URL_ID: 3,
	PANEL_REQ_HEADERS: 4,
	PANEL_REQ_BODY: 5,
	PANEL_RES: 6,
}

const (
	PANEL_UNFOCUS_COLOR string = "240"
	PANEL_FOCUS_COLOR string = "FFF"
)

//mother of panel i.e: the interface that every panel type should follow
type BasePanel interface {
	Init() tea.Cmd
	Update(tea.Msg) tea.Cmd
	View() string
}

type BorderConfig struct {
	Enabled bool
	Width int
	Height int
	Border  lipgloss.Border // if zero, use NormalBorder
	Color   lipgloss.Color
	Top    bool
	Bottom bool
	Left   bool
	Right  bool
}

type MarginConfig struct {
	Top    int
	Bottom int 
	Left   int  
	Right  int 
}

const (
	BORDER_UP    = "BORDER_UP"
	BORDER_DOWN  = "BORDER_DOWN"
	BORDER_LEFT  = "BORDER_LEFT"
	BORDER_RIGHT = "BORDER_RIGHT"
)
