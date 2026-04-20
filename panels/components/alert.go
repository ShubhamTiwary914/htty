package components

import (
	types "htty/types"
	global "htty/globals"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "charm.land/lipgloss/v2"
)

type AlertLevel int

const (
	AlertInfo AlertLevel = 0
	AlertSuccess AlertLevel = 1
	AlertWarn AlertLevel = 2
	AlertError AlertLevel = 3
)

var alertLevelColors = map[AlertLevel]string{
	AlertInfo:    "#5f9ea0",
	AlertSuccess: "#3cb371",
	AlertWarn:    "#daa520",
	AlertError:   "#cd5c5c",
}

var alertLevelLabels = map[AlertLevel]string{
	AlertInfo:    " info ",
	AlertSuccess: " success ",
	AlertWarn:    " warn ",
	AlertError:   " error ",
}

type alertDismissMsg struct{}


type AlertPane struct {
	Message string
	Level   AlertLevel
	TTL     time.Duration
	EndKey  string
	Active bool
	Dimensions types.PaneGeometry	
} 

func (alert *AlertPane) Init() tea.Cmd {
	return nil;
}

func (alert *AlertPane) Update(msg tea.Msg) tea.Cmd {
	if !alert.Active {
		return nil
	}
	switch msg := msg.(type) {
	//event after duration > ttl or endkey
	case alertDismissMsg:
		alert.Dismiss()
	case tea.KeyMsg:
		if alert.EndKey != "" && msg.String() == alert.EndKey {
			alert.Dismiss()
		}
	}
	return nil
}

func (alert AlertPane) View() string {
	if !alert.Active {
		return ""
	}
	color := alertLevelColors[alert.Level]
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true, true, true, true).
		BorderForeground(lipgloss.Color(color)).
		Background(lipgloss.Color(global.Config.Common.Background_color)).
		Width(alert.Dimensions.Width).
		Padding(0, 1).
		Render("press <" + alert.EndKey + "> to close window\n\n" +  alert.Message)
	return box
}

//make lipgloss compositor window from dimensions (iff active)
func (alert AlertPane) ViewAsLayer() *lipgloss.Layer {
	if !alert.Active {
		return nil
	}
	return lipgloss.NewLayer(alert.View()).
		X(alert.Dimensions.X).
		Y(alert.Dimensions.Y).
		Z(3)
}


//activates the alert with a message and level, and kicks off the TTL tick if set.
func (alert *AlertPane) Show(msg string, level AlertLevel) tea.Cmd {
	alert.Message = msg
	alert.Level = level
	alert.Active = true
	if alert.TTL > 0 {
		return tea.Tick(alert.TTL, func(t time.Time) tea.Msg {
			return alertDismissMsg{}
		})
	}
	return nil
}
func (alert *AlertPane) Dismiss() {
	alert.Active = false
	alert.Message = ""
}


//send event to [global.EVENT_ALERTPANE] so app's handler can render the alert dialog
func MakeAlert(message string, seconds int, dimensions types.PaneGeometry, level AlertLevel){
	global.StateBus.Publish(global.EVENT_ALERTPANE, AlertPane{
		Message: message, 
		EndKey: "d", 
		TTL: time.Duration(seconds) * time.Second,
		Dimensions: dimensions,
		Level: level,
	})
}

