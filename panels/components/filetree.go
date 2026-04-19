package components

import (
	"os"
	global "htty/globals"
	types "htty/types"
	utils "htty/utils"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type FileTree struct {
	Width, Height int
	Picker        filepicker.Model
	PanelID       string
	PanelTitle    string
	Border        types.BorderConfig
	Margin        types.MarginConfig
	StatusOptions []string
	err           error
}

func (ft *FileTree) Init() tea.Cmd {
	picker := filepicker.New()
	cwd, err := os.Getwd()
	if err != nil {
		cwd, _ = os.UserHomeDir()
	}
	picker.CurrentDirectory = cwd
	picker.ShowPermissions = false
	picker.ShowSize = false
	picker.FileAllowed = true
	picker.DirAllowed = false 
	picker.AllowedTypes = append(picker.AllowedTypes, global.STATE_FILETYPE);

	ft.Picker = picker
	return ft.Picker.Init()
}

func (ft *FileTree) Update(msg tea.Msg) tea.Cmd {
    var cmd tea.Cmd
    focused := global.CurrentPanelID == global.PANEL_FOCUS_IDS[ft.PanelID]
    if focused {
        utils.SetStatusLineOptions(ft.StatusOptions)
    }

    switch msg.(type) {
    case tea.KeyMsg:
        if focused {
            ft.Picker, cmd = ft.Picker.Update(msg)
        }
    default:
        ft.Picker, cmd = ft.Picker.Update(msg)
    }
    return cmd
}



func (ft FileTree) View() string {
	ft.Border.Color = utils.GetPanelFocusColor(ft.PanelID)
	ft.Picker.SetHeight(ft.Height - 4)
	return ft.Picker.View() 
}

func (ft *FileTree) SetSize(width, height int) {
	ft.Width = width
	ft.Height = height

	w := width - 10
	ft.Picker.Styles.File      = ft.Picker.Styles.File.Width(w)
	ft.Picker.Styles.Directory = ft.Picker.Styles.Directory.Width(w)
	ft.Picker.Styles.Selected  = ft.Picker.Styles.Selected.Width(w)
}
