// Package cli provides a Bubble Tea model with a simple list component.
package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Define the application states
type state int

const (
	stateProjectType state = iota
	stateWorkflowName
	stateFinalize
)

// Define the model struct
type model struct {
	state     state
	list      list.Model
	textInput textinput.Model
}

// NewModel initializes the model with a list and text input component
func NewModel() model {
	items := []list.Item{
		item("Go"),
		item("Python"),
		item("Node.js"),
		item("None of the Above (default Yaml config)"),
	}

	// Initialize the list with default styles
	l := list.New(items, list.NewDefaultDelegate(), 50, 15)
	l.Title = "Select the programming language:"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	// Initialize the text input
	ti := textinput.New()
	ti.Placeholder = "Enter a name for this workflow"
	ti.CharLimit = 64
	ti.Width = 40

	return model{
		state:     stateProjectType,
		list:      l,
		textInput: ti,
	}
}

// Init initializes the program
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateProjectType:
		m.list, cmd = m.list.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				selectedItem := m.list.SelectedItem()
				if selectedItem != nil {
					m.state = stateWorkflowName
					m.textInput.Focus()
				}
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}
	case stateWorkflowName:
		m.textInput, cmd = m.textInput.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.state = stateFinalize
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}
	case stateFinalize:
		return m, tea.Quit
	}
	return m, cmd
}

// View renders the UI based on the current state
func (m model) View() string {
	switch m.state {
	case stateProjectType:
		return m.list.View()
	case stateWorkflowName:
		return fmt.Sprintf(
			"Enter a name for this workflow:\n\n%s\n\n(Press Enter to continue)",
			m.textInput.View(),
		)
	case stateFinalize:
		return "Workflow setup completed! Press Ctrl+C to exit."
	default:
		return "An unexpected error occurred."
	}
}

// item implements the list.Item interface
type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }
