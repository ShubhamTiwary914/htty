package components

import (
	"slices"
	global "htty/globals"
	types "htty/types"
	utils "htty/utils"

	"strings"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)


var headerDelimiters = []rune{':', ' ', '\n'}

type HeaderPane struct {
	Input           textarea.Model
	CharLimit       int
	PanelTitle      string
	PanelID         string
	Placeholder     string
	Border          types.BorderConfig
	StatusOptions   []string
	OptionsStore    []string
	OptionsFilePath string
	OptionBuffer    []string
	
	SuggestionsLimit int 
	selectIndex int
	currentWord string  // word under cursor, recomputed on each change
	cursorLine  int     // Input.Line() snapshot used for overlay Y
	cursorCol   int     // word-start col used for overlay X

	Dimensions types.PaneGeometry
	PaneCfg    types.HttyPanel
}


func (h *HeaderPane) Init() tea.Cmd {
	h.CharLimit = 2052635
	h.Border = types.BorderConfig{Top: true, Bottom: true, Left: true, Right: true}
	h.StatusOptions = []string{}

	input := textarea.New()
	input.Placeholder = h.Placeholder
	input.ShowLineNumbers = false
	input.CharLimit = h.CharLimit
	input.Prompt = ""
	h.Input = input

	var err error
	h.OptionsStore, err = utils.ReadTextLines_intoList(h.OptionsFilePath)
	if err != nil {
		utils.Errorf("header options file %s could not be opened: %v", h.OptionsFilePath, err)
	} else {
		utils.Debugf("success opening header options file %s", h.OptionsFilePath)
	}
	return nil
}


func (h *HeaderPane) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	focused := global.CurrentPanelID == global.PANEL_FOCUS_IDS[h.PanelID]

	if focused {
		h.Input.Focus()
		h.Input.Prompt = "│ "
		utils.SetStatusLineOptions(h.StatusOptions)
	} else {
		h.Input.Blur()
		h.Input.Prompt = ""
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case global.Config.Key.Textoptions_prev:
			if focused && len(h.OptionBuffer) > 0 {
				h.selectIndex--
				if h.selectIndex < 0 {
					h.selectIndex = len(h.OptionBuffer) - 1
				}
				return cmd
			}
		case global.Config.Key.Textoptions_next:
			if focused && len(h.OptionBuffer) > 0 {
				h.selectIndex++
				if h.selectIndex >= len(h.OptionBuffer) {
					h.selectIndex = 0
				}
				return cmd
			}
		case global.Config.Key.CompleteText:
			if focused && len(h.OptionBuffer) >= 1 {
				h.__acceptCompletion()
				h.clearOptions()
				return cmd
			}
		}
	}

	oldValue := h.Input.Value()
	h.Input, cmd = h.Input.Update(msg)
	newValue := h.Input.Value()

	if newValue == "" {
		h.clearOptions()
		return cmd
	}
	if newValue != oldValue {
		h.__recomputeWord()
	}
	return cmd
}


func (h *HeaderPane) View() string {
	baseView, _ := h.ViewWithOptions(false)
	style := lipgloss.NewStyle().Background(lipgloss.Color(global.Config.Common.Background_color))
	return style.Render(baseView)
}


func (h *HeaderPane) SetSize() {
	h.Input.SetWidth(h.Dimensions.Width)
	h.Input.SetHeight(h.Dimensions.Height)
}


func (h *HeaderPane) ViewWithOptions(withLayer bool) (baseView string, optionsLayer *lipgloss.Layer) {
	h.Border.Color = utils.GetPanelFocusColor(h.PanelID)
	style := utils.SetBorder(h.Border).
		BorderForeground(lipgloss.Color(h.Border.Color)).
		Background(lipgloss.Color(global.Config.Common.Background_color))

	baseView = utils.SetBorderStyle_WithLabelTop(style, h.Input.View(), h.Border,
		utils.GetPanelTitleLabel(h.PanelTitle, global.PANEL_FOCUS_IDS[h.PanelID]),
	)

	if !withLayer {
		return baseView, nil
	}

	focused := global.CurrentPanelID == global.PANEL_FOCUS_IDS[h.PanelID]
	if !focused || len(h.OptionBuffer) == 0 {
		return baseView, nil
	}

	optionsStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true, true, true).
		BorderForeground(lipgloss.Color(global.Config.Common.Focus_border_color)).
		Background(lipgloss.Color(global.Config.Common.Background_color)).
		Width(h.Dimensions.Width / 2)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(global.Config.Common.Textoptions_selection_color))

	var lines []string
	for i, opt := range h.OptionBuffer {
		if i == h.selectIndex {
			lines = append(lines, selectedStyle.Render(opt))
		} else {
			lines = append(lines, opt)
		}
	}

	newContent := strings.Join(lines, "\n")
	optionsBox := optionsStyle.Render(newContent)

	// position below the cursor line, offset X to word-start column
	// +1 for top border, +1 to place below the cursor line
	overlayX := h.Dimensions.X + 1 + h.cursorCol
	overlayY := h.Dimensions.Y + 1 + h.cursorLine + 1

	optionsLayer = lipgloss.NewLayer(optionsBox).X(overlayX).Y(overlayY).Z(2)
	return baseView, optionsLayer
}


func (h *HeaderPane) SetValue(value string) { h.Input.SetValue(value) }
func (h *HeaderPane) GetValue() string      { return h.Input.Value() }


// __recomputeWord extracts the word under the cursor and updates
// OptionBuffer, cursorLine, and cursorCol for overlay positioning.
func (h *HeaderPane) __recomputeWord() {
	lineInfo := h.Input.LineInfo()
	col := lineInfo.CharOffset
	h.cursorLine = h.Input.Line()

	// get just the current line up to cursor
	lines := strings.Split(h.Input.Value(), "\n")
	if h.cursorLine >= len(lines) {
		h.clearOptions()
		return
	}
	lineUpToCursor := []rune(lines[h.cursorLine])
	if col > len(lineUpToCursor) {
		col = len(lineUpToCursor)
	}
	lineUpToCursor = lineUpToCursor[:col]

	// scan backwards to find word start
	wordStart := len(lineUpToCursor)
	for wordStart > 0 && !isHeaderDelimiter(lineUpToCursor[wordStart-1]) {
		wordStart--
	}

	word := string(lineUpToCursor[wordStart:])
	h.cursorCol = wordStart  // overlay X anchored at word start
	h.currentWord = word

	if word == "" {
		h.clearOptions()
		return
	}

	var err error
	h.OptionBuffer, err = utils.PrefixClosestSearch_withOptions(word, h.OptionsStore)
	if len(h.OptionBuffer) > h.SuggestionsLimit {
		h.OptionBuffer = h.OptionBuffer[:h.SuggestionsLimit]
	}
	if err != nil {
		panic(err)
	}
	utils.Debugf("res: %v, from: %v", h.OptionBuffer, h.OptionsStore)
	h.selectIndex = 0
}


// __acceptCompletion replaces the current word in the textarea with the
// selected completion. If the cursor is before any ':' on the current line
// (i.e. we're completing a header key), ": " is appended automatically.
func (h *HeaderPane) __acceptCompletion() {
	if len(h.OptionBuffer) == 0 {
		return
	}
	completion := h.OptionBuffer[h.selectIndex]

	lines := strings.Split(h.Input.Value(), "\n")
	if h.cursorLine >= len(lines) {
		return
	}

	currentLine := lines[h.cursorLine]
	lineInfo := h.Input.LineInfo()
	col := min(lineInfo.CharOffset, len([]rune(currentLine)))

	// if key, put a ":" at end, if val then break line to next header
	beforeCursor := string([]rune(currentLine)[:col])
	isKey := !strings.Contains(beforeCursor, ":")
	if isKey {
		completion = completion + ": "
	} else {
		completion = completion + "\n"
	}

	// replace [wordStart:col] with the completion
	lineRunes := []rune(currentLine)
	newLine := string(lineRunes[:h.cursorCol]) + completion + string(lineRunes[col:])
	lines[h.cursorLine] = newLine

	h.Input.SetValue(strings.Join(lines, "\n"))
}


func (h *HeaderPane) clearOptions() {
	h.selectIndex = 0
	h.OptionBuffer = nil
	h.currentWord = ""
}

func isHeaderDelimiter(r rune) bool {
	return slices.Contains(headerDelimiters, r)
}

