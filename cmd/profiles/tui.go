package profiles

import (
	"fmt"
	"os"

	common "github.com/hectorruiz-it/grabber/cmd"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func profileTui(basic bool, ssh bool, token bool) model {
	p := tea.NewProgram(initialModel(basic, ssh, token))

	finalModel, err := p.Run()
	common.CheckAndReturnError(err)
	m := finalModel.(model)

	return m
}

type (
	errMsg error
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#fed858"))
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type model struct {
	inputs      []textinput.Model
	focused     int
	err         error
	profileType string
}

func initialModel(basic bool, ssh bool, token bool) model {
	var profileType string

	var m model
	// var inputs []textinput.Model = make([]textinput.Model, 6)

	switch {
	case basic:
		var inputs []textinput.Model = make([]textinput.Model, 3)
		for i := range inputs {
			if i == 0 {
				inputs[i] = textinput.New()
				inputs[i].Placeholder = "git"
				inputs[i].Focus()
				inputs[i].CharLimit = 50
				inputs[i].Width = 60
				inputs[i].Prompt = ""
			}
			if i == 1 {
				inputs[i] = textinput.New()
				inputs[i].Placeholder = "XXXXXXXXXXXX"
				inputs[i].CharLimit = 60
				inputs[i].Width = 60
				inputs[i].Prompt = ""
				inputs[i].EchoMode = textinput.EchoPassword
				inputs[i].EchoCharacter = '•'
			}
			if i == 2 {
				inputs[i] = textinput.New()
				inputs[i].Placeholder = ""
				inputs[i].CharLimit = 0
				inputs[i].Width = 0
				inputs[i].Prompt = "[ Apply changes ]"
				inputs[i].PromptStyle = continueStyle
			}
		}
		profileType = "basic"

		m = model{
			inputs:      inputs,
			focused:     0,
			err:         nil,
			profileType: profileType,
		}

	case ssh:
		var inputs []textinput.Model = make([]textinput.Model, 4)
		for i := range inputs {
			if i == 0 {
				inputs[i] = textinput.New()
				inputs[i].Placeholder = "/home/grabber/.ssh/github"
				inputs[i].Focus()
				inputs[i].CharLimit = 50
				inputs[i].Width = 60
				inputs[i].Prompt = ""
			}
			if i == 1 {
				inputs[i] = textinput.New()
				inputs[i].Placeholder = "/home/grabber/.ssh/github.pub"
				inputs[i].CharLimit = 50
				inputs[i].Width = 60
				inputs[i].Prompt = ""
			}
			if i == 2 {
				inputs[i] = textinput.New()
				inputs[i].Placeholder = "XXXXXXXXXXXX"
				inputs[i].CharLimit = 60
				inputs[i].Width = 60
				inputs[i].Prompt = ""
				inputs[i].EchoMode = textinput.EchoPassword
				inputs[i].EchoCharacter = '•'
			}
			if i == 3 {
				inputs[i] = textinput.New()
				inputs[i].Placeholder = ""
				inputs[i].CharLimit = 0
				inputs[i].Width = 0
				inputs[i].Prompt = "[ Apply changes ]"
				inputs[i].PromptStyle = continueStyle
			}
		}
		profileType = "ssh"

		m = model{
			inputs:      inputs,
			focused:     0,
			err:         nil,
			profileType: profileType,
		}

	case token:
		var inputs []textinput.Model = make([]textinput.Model, 3)
		for i := range inputs {
			if i == 0 {
				inputs[i] = textinput.New()
				inputs[i].Placeholder = "git"
				inputs[i].Focus()
				inputs[i].CharLimit = 50
				inputs[i].Width = 60
				inputs[i].Prompt = ""
			}
			if i == 1 {
				inputs[i] = textinput.New()
				inputs[i].Placeholder = "XXXXXXXXXXXX"
				inputs[i].CharLimit = 60
				inputs[i].Width = 60
				inputs[i].Prompt = ""
				inputs[i].EchoMode = textinput.EchoPassword
				inputs[i].EchoCharacter = '•'
			}
			if i == 2 {
				inputs[i] = textinput.New()
				inputs[i].Placeholder = ""
				inputs[i].CharLimit = 0
				inputs[i].Width = 0
				inputs[i].Prompt = "[ Apply changes ]"
				inputs[i].PromptStyle = continueStyle
			}
		}

		profileType = "token"

		m = model{
			inputs:      inputs,
			focused:     0,
			err:         nil,
			profileType: profileType,
		}

	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				err, errIndex := m.ValidateInputs()
				if err != nil {
					m.err = err
					m.focused = errIndex
					cmds := make([]tea.Cmd, len(m.inputs))
					for i := 0; i <= len(m.inputs)-1; i++ {
						if i == m.focused {
							// Set focused state
							cmds[i] = m.inputs[i].Focus()
							continue
						}
						// Remove focused state
						m.inputs[i].Blur()
					}
					return m, tea.Batch(cmds...)
				}

				// storeOnKeychain(m.inputs[passphrase].Value())
				return m, tea.Quit
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			os.Exit(1)
		case tea.KeyShiftTab, tea.KeyCtrlP, tea.KeyUp:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN, tea.KeyDown:
			m.nextInput()
		}
		if m.focused == len(m.inputs)-1 {
			m.inputs[m.focused].Cursor.SetMode(cursor.CursorHide)
			m.inputs[m.focused].PromptStyle = inputStyle
		} else {
			m.inputs[m.focused].PromptStyle = continueStyle
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string

	switch m.profileType {
	case "ssh":
		s = fmt.Sprintf(` %s
 %s

 %s
 %s

 %s
 %s

 %s
`,
			inputStyle.Width(50).Render("Private Key Path"),
			m.inputs[0].View(),
			inputStyle.Width(50).Render("Public Key Path"),
			m.inputs[1].View(),
			inputStyle.Width(50).Render("Passprhase"),
			m.inputs[2].View(),
			m.inputs[3].View(),
		) + "\n"
	case "basic":
		s = fmt.Sprintf(` %s
 %s

 %s
 %s

 %s
`,
			inputStyle.Width(50).Render("Username"),
			m.inputs[0].View(),
			inputStyle.Width(50).Render("Passphrase"),
			m.inputs[1].View(),
			m.inputs[2].View(),
		) + "\n"
	case "token":
		s = fmt.Sprintf(` %s
 %s

 %s
 %s

 %s
`,
			inputStyle.Width(50).Render("Username"),
			m.inputs[0].View(),
			inputStyle.Width(50).Render("Token"),
			m.inputs[1].View(),
			m.inputs[2].View(),
		) + "\n"
	}

	if m.err != nil {
		return fmt.Sprintf("%s %s", s, m.err)
	}

	return s
}

// nextInput focuses the next input field
func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *model) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func (m *model) ValidateInputs() (error, int) {
	// if m.inputs[0].Value() == "" {
	// 	return errors.New("grabber: private key path is required"), 0
	// }

	// if err := checkPathAbsoluteExists(m.inputs[0].Value(), "ssh key"); err != nil {
	// 	return err, 0
	// }

	// if m.inputs[1].Value() == "" {
	// 	return errors.New("grabber: public key path is required"), 1
	// }

	// if err := checkPathAbsoluteExists(m.inputs[1].Value(), "ssh public key"); err != nil {
	// 	return err, 1
	// }

	// if m.inputs[2].Value() == "" {
	// 	return errors.New("grabber: passhprase is required"), 2
	// }
	return nil, -1
}
