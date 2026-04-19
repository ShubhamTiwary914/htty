package components 

import (
	global "htty/globals"
	types "htty/types"
	utils "htty/utils"

	"strings"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type TextOptions struct {
	Width, Height int
	Input         textarea.Model
	CharLimit     int
	PanelTitle    string
	PanelID       string
	Placeholder   string
	Showline      bool
	Border        types.BorderConfig
	OptionBuffer   []string 
	OptionsStore   []string   
	OptionsFilePath string
	AllowSaveInput    bool //whether to save inputs for later options once api call made
	StatusOptions []string

	selectIndex   int

	Dimensions    types.PaneGeometry
	PaneCfg       types.HttyPanel

}

func (text *TextOptions) Init() tea.Cmd {
	text.CharLimit = 2052635
	text.Border = types.BorderConfig{Top: true, Bottom: true, Left: true, Right: true}
	text.Showline = false
	text.StatusOptions = []string{}

	var input textarea.Model = textarea.New()
	input.Placeholder = text.Placeholder
	input.ShowLineNumbers = text.Showline
	input.CharLimit = text.CharLimit
	input.Prompt = ""
	text.Input = input

	var err error
	text.OptionsStore, err = utils.ReadTextLines_intoList(text.OptionsFilePath)
	if err != nil { 
		utils.Errorf("file %s could not be opened for options store, err: %v", text.OptionsFilePath, err); 
	} else {  utils.Debugf("file %s successfully opened for options store", text.OptionsFilePath); }
	return nil
}

func (text *TextOptions) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	focused := global.CurrentPanelID == global.PANEL_FOCUS_IDS[text.PanelID]
	
	if focused {
		utils.SetStatusLineOptions(text.StatusOptions)
		text.Input.Focus()
		text.Input.Prompt = "│ "
	} else {
		text.Input.Blur()
		text.Input.Prompt = ""
	}

	switch msg := msg.(type) {
	//auto resize to window dimensions
	case tea.KeyMsg:
		switch msg.String(){
		case global.Config.Key.Textoptions_prev:
			text.selectIndex--
			if text.selectIndex < 0 {
				text.selectIndex = len(text.OptionBuffer) - 1  // Wrap to end
			}
			return cmd
		case global.Config.Key.Textoptions_next:
			text.selectIndex++
			if text.selectIndex >= len(text.OptionBuffer) {
				text.selectIndex = 0  // Wrap to beginning
			}
			return cmd
		case global.Config.Key.CompleteText:
			if focused && len(text.OptionBuffer)>=1{
				text.Input.SetValue(text.OptionBuffer[text.selectIndex])
				text.ClearOptions()
				return cmd
			}
		case global.Config.Key.Sendapicall:
			if text.AllowSaveInput{
				err := utils.Insert_SortedFile(text.Input.Value(), text.OptionsFilePath)
				utils.Debugf("saving %s into file %s", text.Input.Value(), text.OptionsFilePath)
				if err != nil {
					utils.Errorf("insertion of %s into file(%s) failed!", text.Input.Value(), text.OptionsFilePath)
				}
			}
		}
	}
	
	oldValue := text.Input.Value()
	text.Input, cmd = text.Input.Update(msg)
	if text.Input.Value() == "" {
		text.ClearOptions()
		return cmd;
	}
	if text.Input.Value() != oldValue {
		text.SetOptions(text.Input.Value())
	}
	return cmd
}

// View renders just the text input (for use when called directly)
func (text TextOptions) View() string {
	baseView, _ := text.ViewWithOptions(false)
	style := lipgloss.NewStyle().Background(lipgloss.Color(global.Config.Common.Background_color))
	return style.Render(baseView)
}

// ViewWithOptions returns the base view AND optionally the options layer
// If withLayer is true, returns the options layer for compositor
// If withLayer is false, returns nil layer (used when View() is called directly)
func (text TextOptions) ViewWithOptions(withLayer bool) (baseView string, optionsLayer *lipgloss.Layer) {
	// Render the base text input
	inputStyle := utils.SetBorder(text.Border).BorderForeground(
		lipgloss.Color(utils.GetPanelFocusColor(text.PanelID)),
	).Background(lipgloss.Color(global.Config.Common.Background_color))
	
	text.Border.Color = utils.GetPanelFocusColor(text.PanelID)
	baseView = utils.SetBorderStyle_WithLabelTop(inputStyle, text.Input.View(), text.Border, 
		utils.GetPanelTitleLabel(text.PanelTitle, global.PANEL_FOCUS_IDS[text.PanelID]),		
	)
	// if not requesting layer, return early
	if !withLayer {
		return baseView, nil
	}
	
	focused := global.CurrentPanelID == global.PANEL_FOCUS_IDS[text.PanelID]
	
	// If not focused or no options, return nil layer
	if !focused || len(text.OptionBuffer) == 0 {
		return baseView, nil
	}

	// create options overlay
	optionsStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true, true, true).
		BorderForeground(lipgloss.Color(global.Config.Common.Focus_border_color)).
		Background(lipgloss.Color(global.Config.Common.Background_color)).
		Width(text.Width)

	// style for selected item
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(global.Config.Common.Textoptions_selection_color))

	// build options text with selected item highlighted
	var optionsLines []string
	for i, option := range text.OptionBuffer {
		if i == text.selectIndex{

			optionsLines = append(optionsLines, selectedStyle.Render(option))
		} else {
			optionsLines = append(optionsLines, option)
		}
	}
	optionsText := strings.Join(optionsLines, "\n")
	optionsBox := optionsStyle.Render(optionsText)
	// return the options layer with Z=2 (above base panels)
	optionsLayer = lipgloss.NewLayer(optionsBox).Z(2)
	return baseView, optionsLayer
}

func (text *TextOptions) SetSize() {
	text.Input.SetWidth(text.Dimensions.Width)
	text.Input.SetHeight(text.Dimensions.Height)
}

//change the Options[] with new suggestions
func (text *TextOptions) SetOptions(inputstr string) {
	text.OptionBuffer, _ = utils.PrefixClosestSearch_withOptions(inputstr, text.OptionsStore)	
	text.selectIndex = 0
}

func (text *TextOptions) ClearOptions(){
	text.selectIndex = 0
	text.OptionBuffer = nil
}
