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
	stateCloudProvidier
)

// Define the model struct
type model struct {
	state          state
	supportedLangs list.Model
	supprotedCloud list.Model
	textInput      textinput.Model
	workflowName   string
	langauge       string                 //might chnage to a map so that it maps a string -> the github action for that programming language
	cloud          string                 //also above
	secrets        map[string]interface{} //also above

}

// NewModel initializes the model with a list and text input component
func NewModel() model {
	langs := []list.Item{
		item("Go"),
		item("Python"),
		item("Node.js"),
		item("None of the Above (default Yaml config)"),
	}

	cloud_providers := []list.Item{
		item("AWS"),
		item("Azure"),
		item("GCP"),
		item("None of the Above"),
	}

	// Initialize the list with default styles
	languages := list.New(langs, list.NewDefaultDelegate(), 50, 15)
	languages.Title = "Select the programming language:"
	languages.SetShowStatusBar(false)
	languages.SetShowHelp(false)

	cloud := list.New(cloud_providers, list.NewDefaultDelegate(), 50, 15)
	cloud.Title = "Select a Cloud provider to configure sso credentials:"
	cloud.SetShowStatusBar(false)
	cloud.SetShowHelp(false)

	// Initialize the text input
	ti := textinput.New()
	ti.Placeholder = "Enter a name for this workflow"
	ti.CharLimit = 64
	ti.Width = 40

	return model{
		state:          stateProjectType,
		supportedLangs: languages,
		supprotedCloud: cloud,
		textInput:      ti,
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
		m.supportedLangs, cmd = m.supportedLangs.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				selectedItem := m.supportedLangs.SelectedItem()
				if selectedItem != nil {
					m.state = stateCloudProvidier
					m.textInput.Focus()
					m.langauge = selectedItem.FilterValue()
				}
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}
	case stateCloudProvidier:
		m.supprotedCloud, cmd = m.supprotedCloud.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				selectedCloud := m.supprotedCloud.SelectedItem()
				if selectedCloud != nil {
					m.state = stateWorkflowName
					m.textInput.Focus()
					m.cloud = selectedCloud.FilterValue()
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
			m.workflowName = m.textInput.Value()
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
		return m.supportedLangs.View()

	case stateCloudProvidier:
		return m.supprotedCloud.View()
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
