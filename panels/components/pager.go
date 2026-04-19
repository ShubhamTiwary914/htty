package components 

import (
	types "htty/types"

	global "htty/globals"

	utils "htty/utils"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)


type PagerPane struct {
	Width, Height int
	Viewport      viewport.Model
	PanelTitle    string
	PanelID       string
	Content       string
	StatusOptions []string
	Border        types.BorderConfig
	Margin        types.MarginConfig
	ready         bool
}

func (pager *PagerPane) Init() tea.Cmd {
	pager.Viewport = viewport.New(
		pager.Width,
		pager.Height,
	)
	pager.Viewport.SetContent(pager.Content)
	pager.ready = true
	return nil
}

func (pager *PagerPane) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	focused := global.CurrentPanelID == global.PANEL_FOCUS_IDS[pager.PanelID]
	if focused {
		utils.SetStatusLineOptions(pager.StatusOptions)
	}
	switch msg.(type) {
	case tea.KeyMsg, tea.MouseMsg:
		if focused {
			pager.Viewport, cmd = pager.Viewport.Update(msg)
		}
	default:
		pager.Viewport, cmd = pager.Viewport.Update(msg)
	}
	return cmd
}

func (pager *PagerPane) View() string {
	pager.Border.Color = utils.GetPanelFocusColor(pager.PanelID)
	style := utils.SetBorder(pager.Border).
		BorderForeground(lipgloss.Color(pager.Border.Color)).
		Background(lipgloss.Color(global.Config.Common.Background_color))
	return utils.SetBorderStyle_WithLabelTop(style, pager.Viewport.View(), pager.Border,
		utils.GetPanelTitleLabel(pager.PanelTitle, global.PANEL_FOCUS_IDS[pager.PanelID]),
	)
}


func (pager *PagerPane) SetSize(width, height int) {
	pager.Width = width
	pager.Height = height
	pager.Viewport.Width = width
	pager.Viewport.Height = height
}

func (pager *PagerPane) SetContent(content string) {
	pager.Content = content
	pager.Viewport.SetContent(content)
	pager.Viewport.GotoTop()
}

func (pager *PagerPane) ScrollPercent() float64 {
	return pager.Viewport.ScrollPercent()
}

