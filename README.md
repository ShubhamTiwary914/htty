
# HttY - TUI toolkit for HTTP calls over curl 

### Usage for local run
- Install go (recommened: 1.26.0)

Install packages:
```bash
go install
```

Run local:
```bash
make dev    #requires make, which most OS comes with inbuilt
```

For full supported commands list:
```bash
❯_ make help
HttY makefile guide (source: Makefile)
Usage: make [Target]

[Targets]
dev         run htty local in dev/debug mode
build       build htty executable
logwatch    follow debug.log for viewing live logs
```

---

### Contributing Guidelines

> will be moved to a seperate "CONTRIBUTING.md", this is it for now

1. For any large feature changes, please update this readme section "Stuff to add" as marked
2. Try to avoid changes to go.mod or go.sum files (i.e: adding new packages) since they cause a lot of conflicts


Structure for the UI side is:
```
App panel --> Main Panel -> sub panels
          --> Side Panel -> sub panels
```
i.e: everything is composed up of panels


Main libraries that you'll need context for are:
- [Bubble tea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)

**About tea**
Tea is a state management tool (like react), it updates and renders these panels on changes (from tea's POV, models == panels), and each model has:
- Init()    - when model starts
- Update()  - handles incoming events and updates the model accordingly.
- View()    - gets triggered when model changes (actual UI changes here)

example:
```go
type MainPane struct {
}

func (main MainPane) Init() tea.Cmd {
	return nil;
}

func (main MainPane) Update(msg tea.Msg) (MainPane, tea.Cmd) {
	return main, nil;
}

func (m MainPane) View() string {
	return "main panel"
}
```

---


### Stuff to add

- [ ] v0.1 (working tool)
  - [ ] main wrapper utility for making HTTP requests with curl
  - [ ] keybinds event handler (hardcoded keys) to navigate between panes/subpanes
  - [ ] file parser to store HttY request (body,headers,..) & file picker on left  [navigation]
      
- [ ] v0.1 addons (additional features)
  - [ ] central config for styling (that user can control --> ideally somewhere like ~/.config/htty/style.toml)
  - [ ] keybinds config (similar to style, maybe somewhere like ~/.config/htty/keybinds.toml)
  - [ ] navigate quickly using pane IDs (like main=1, side=2, ...)

> feel free to add here (cuz I can't remember shit right now)


**About lipgloss**
Lipgloss is simply the UI management, tea is left to deal with model state changes so actions like keypress, etc can trigger Update() -> View () -> which calls lipgloss for styling changes

