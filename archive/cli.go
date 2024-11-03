// Package cli provides a Bubble Tea model with a simple list component.
/*
Code out of date by 11/3/2024
This code previously served as the entire cli file

might be useful in case of looking for creating a bubbletea cli in one file if it is small

reason for archival: file was going to become too large and ultimately did not make sense having one
file as the entire cli
refactored into cli.go -> init.go, model.go, update.go, view.go


*/

package archive

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Define the application states
type state int

const (
	stateWorkflowName state = iota
	stateSchedule
	stateCronFrequency
	stateCloudProvider
	stateCustomCron
	stateComplete
)

// Define the model struct
type model struct {
	state          state
	supportedSched list.Model
	supportedCloud list.Model
	cronFrequency  list.Model
	textInput      textinput.Model
	workflowName   string
	schedule       string
	cloud          string
	customCron     string
}

// NewModel initializes the model with list and text input components
func NewModel() model {
	// Scheduling options
	schedulers := []list.Item{
		item("On dispatch"),
		item("On Pull"),
		item("On Push"),
		item("Cron Schedule"),
	}

	// Cron frequency options if "Cron Schedule" is selected
	cronOptions := []list.Item{
		item("Once a day"),
		item("Once a week"),
		item("Once a month"),
		item("Once a year"),
		item("Other (Enter custom cron)"),
	}

	// Cloud providers
	cloudProviders := []list.Item{
		item("AWS"),
		item("Azure"),
		item("GCP"),
		item("None of the Above"),
	}

	// Initialize the list components
	schedule := list.New(schedulers, list.NewDefaultDelegate(), 50, 15)
	schedule.Title = "Select a scheduling type:"
	schedule.SetShowStatusBar(false)
	schedule.SetShowHelp(false)

	cloud := list.New(cloudProviders, list.NewDefaultDelegate(), 50, 15)
	cloud.Title = "Select a Cloud provider to configure SSO credentials:"
	cloud.SetShowStatusBar(false)
	cloud.SetShowHelp(false)

	cron := list.New(cronOptions, list.NewDefaultDelegate(), 50, 15)
	cron.Title = "Select Cron Frequency:"
	cron.SetShowStatusBar(false)
	cron.SetShowHelp(false)

	// Initialize the text input
	ti := textinput.New()
	ti.Placeholder = "Enter a name for this workflow"
	ti.CharLimit = 64
	ti.Width = 40

	return model{
		state:          stateWorkflowName,
		supportedSched: schedule,
		cronFrequency:  cron,
		supportedCloud: cloud,
		textInput:      ti,
	}
}

// Init initializes the program and starts text input blinking
func (m model) Init() tea.Cmd {
	// Set initial focus on text input for the workflow name
	m.textInput.Focus()
	return textinput.Blink
}

// Update handles messages and updates the model state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateWorkflowName:
		// Ensure text input is focused
		m.textInput.Focus()
		m.textInput, cmd = m.textInput.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.workflowName = m.textInput.Value()
				m.textInput.Reset()
				m.state = stateSchedule
				return m, textinput.Blink
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}
	case stateSchedule:
		m.supportedSched, cmd = m.supportedSched.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				selectedSchedule := m.supportedSched.SelectedItem()
				if selectedSchedule != nil {
					m.schedule = selectedSchedule.FilterValue()
					if m.schedule == "Cron Schedule" {
						m.state = stateCronFrequency
					} else {
						m.state = stateCloudProvider
					}
				}
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}
	case stateCronFrequency:
		m.cronFrequency, cmd = m.cronFrequency.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				selectedFrequency := m.cronFrequency.SelectedItem()
				if selectedFrequency != nil {
					frequency := selectedFrequency.FilterValue()
					if frequency == "Other (Enter custom cron)" {
						// Prompt user for custom cron input
						m.textInput.Placeholder = "Enter custom cron schedule"
						m.textInput.SetValue("")
						m.textInput.Focus() // Ensure focus is on text input
						m.state = stateCustomCron
						return m, textinput.Blink
					} else {
						// Directly go to cloud provider selection if predefined cron is chosen
						m.schedule = frequency
						m.state = stateCloudProvider
					}
				}
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}
	case stateCloudProvider:
		m.supportedCloud, cmd = m.supportedCloud.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				selectedCloud := m.supportedCloud.SelectedItem()
				if selectedCloud != nil {
					m.cloud = selectedCloud.FilterValue()
					m.state = stateComplete
				}
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}
	case stateCustomCron:
		// Handle custom cron input and then proceed to cloud provider selection
		m.textInput, cmd = m.textInput.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.customCron = m.textInput.Value()
				m.schedule = m.customCron
				m.state = stateCloudProvider // Move to cloud selection after custom cron input
				return m, nil
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}
	case stateComplete:
		// End the program on 'enter' or 'ctrl+c'
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter", "ctrl+c", "q":
				// Quit the program when 'enter', 'ctrl+c', or 'q' is pressed in the final state
				return m, tea.Quit
			}
		}
	}
	return m, cmd
}

// View renders the UI based on the current state
func (m model) View() string {
	switch m.state {
	case stateWorkflowName:
		return fmt.Sprintf("Enter a name for this workflow:\n\n%s\n\n(Press Enter to continue)", m.textInput.View())

	case stateSchedule:
		return m.supportedSched.View()

	case stateCronFrequency:
		return m.cronFrequency.View()

	case stateCloudProvider:
		return m.supportedCloud.View()

	case stateCustomCron:
		return fmt.Sprintf("Enter custom cron schedule:\n\n%s\n\n(Press Enter to continue)", m.textInput.View())

	case stateComplete:
		return "Workflow setup completed! Press Enter or Ctrl+C to exit."
	default:
		return "An unexpected error occurred."
	}
}

// item implements the list.Item interface
type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }
