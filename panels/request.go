package htty

import (
	utils "htty/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RequestPane struct {
	width, height int
	focusIndex int
	initialized bool
	textdrop model
}

func (rq *RequestPane) Init() tea.Cmd {
	utils.Infof("request subpanel initialization")
	rq.textdrop = NewModel()
	return nil;
}

func (rq *RequestPane) Update(msg tea.Msg) (tea.Cmd){	
	var cmds []tea.Cmd
	var cmd tea.Cmd;
	cmd = rq.textdrop.Update(msg) 
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...) 
}

func (rq RequestPane) View() string {
	style := utils.SetBorder(rq.width-5, rq.height-2, lipgloss.RoundedBorder())
	return style.Render(rq.textdrop.View())
}

func (rq *RequestPane) SetSize(w int, h int){
	rq.width = w
	rq.height = h
	rq.textdrop.SetSize(w/11, h/12)
}

