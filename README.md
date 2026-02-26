# HttY
Simple, lightweight Terminal UI for making HTTP API calls

---

## Table of Contents
- [Setup Guide](https://github.com/ShubhamTiwary914/htty/tree/master?tab=readme-ov-file#setup--run-locally)
- [Architecture](https://github.com/ShubhamTiwary914/htty/tree/master?tab=readme-ov-file#architecture)
- [Contributing Guidelines](https://github.com/ShubhamTiwary914/htty/tree/master?tab=readme-ov-file#contributing) 

---

## Setup & run locally

#### Tools required:
- Go (recommended ver: 1.26.0)
- make

Install dependencies:
```bash
go install
```

Run in development mode:
```
make dev
```

Build binary (generates a `htty` binary):
```bash
make build
```

Check full options make has:
```
make help
```

---

## Architecture
Keeping the codebase into two parts:
-  logic/utilitues (like logging,http,file) that can be tested
-  UI, which is handled by packages: `bubble tea` & `lipgloss`
	- [bubble tea](https://github.com/charmbracelet/bubbletea) handles state changes(like react), events(key press, etc)
 	- [lipgloss](https://github.com/charmbracelet/lipgloss) does the actual styling(border, colors, ...)
  	- [bubbles](https://github.com/charmbracelet/bubbles?tab=readme-ov-file#file-picker) [addon]: has ready TUI components (using tea+lipgloss) like file picker, list, ...
  
So we break everything down to 'panels' (can be nested) to work modularly as possible, like:
```bash
App panel
├── Main Panel
│ └── Sub panels
└── Side Panel
└── Sub panels
```

And each panel implements three methods:
| Method   | Purpose                                                                                      |
| -------- | -------------------------------------------------------------------------------------------- |
| Init()   | Runs once when the model starts, for initial commands or setup.                              |
| Update() | Handles incoming messages and state transitions, returns updated model                       |
| View()   | Renders the UI based on current state, called after state changes.                           |

and "state" is a struct that is the "model" of that panel, example:
```go
type MainPane struct {
	somevar int;
}

func (m MainPane) Init() tea.Cmd {
	return nil
}

func (m MainPane) Update(msg tea.Msg) (MainPane, tea.Cmd) {
	m.var = somevalue
	return m, nil
}

func (m MainPane) View() string {
	return "main panel"
    //lipgloss part of styling comes here
}
```

---

## Contributing 

#### Bug reporrs  
Currently `htty` is under development and it will have a lot of bugs, that we will be fixing as they come up.

Feel free to let us know if you encounter any bug, it would help get the tool better. 
To report a bug just open an [issue](https://github.com/ShubhamTiwary914/htty/issues) describing the bug.

#### Feature additions
1. Making an issue for a feature initially with some clear description
2. Fork the repo, making changes & push your branch and open a pull request against the main branch. (ex: `htty/main <- your-repo/feat/something`)
3. Provide a clear description of the change, steps to verify that's successful and reference to the issue: (ex: `closes #10, <some short comment if needed>` or `referencing #100`)


#### Debugging & Testing

> TUI blocks stdout, so logs are piped to a `,/htty.log` file (when LOGLEVEL env is set while running htty)

View live logs with:
> 
```bash
make logwatch
```

> Tests are in `./tests/` folder which ideally, all core utilities should have tests for, run tests with:
```bash
make test
``` 

for reference: [Guide for writing tests in golang](https://go.dev/doc/tutorial/add-a-test)

---

### LICENSE

This project is licensed under the MIT License.

See [LICENSE](https://github.com/ShubhamTiwary914/htty/blob/master/LICENSE) for full details.
