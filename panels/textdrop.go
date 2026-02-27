package htty

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	utils "htty/utils"
	"github.com/charmbracelet/lipgloss"
)
var options = []string{"foo", "bar", "baz"}

type model struct {
	width, height int
	ta textarea.Model
}

func NewModel() model {
	ta := textarea.New()
	ta.Placeholder = "Method"
	ta.Focus()
	ta.ShowLineNumbers = false
	ta.SetHeight(1)
	ta.CharLimit = 10 
	ta.Prompt = ""

	return model{
		ta: ta,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.ta, cmd = m.ta.Update(msg)
	return cmd
}

func (m model) View() string {
	BorderBottom := lipgloss.Border{
		Bottom: "─", 
		Top:    " ", 
		Left:   " ",
		Right:  " ",
	}
	style := utils.SetBorder(m.width, m.height, BorderBottom) 
	return style.Render(m.ta.View())
}

func (m *model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.ta.SetWidth(width)
	m.ta.SetHeight(height)
}
