// init.go
package cli

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// NewModel initializes the model with list and text input components
func NewModel() model {
	// Scheduling options
	schedulers := []list.Item{
		item("On dispatch"),
		item("On Pull"),
		item("On Push"),
		item("Cron Schedule"),
	}

	// Supported programming languages
	languages := []list.Item{
		item("Go"),
		item("Python"),
		item("Node.js"),
		item("Java"),
	}

	lang := list.New(languages, list.NewDefaultDelegate(), 50, 15)
	lang.Title = "Select a programming language:"
	lang.SetShowStatusBar(false)
	lang.SetShowHelp(false)

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

	// Initialize list components
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

	// Initialize text inputs
	ti := textinput.New()
	ti.Placeholder = "Enter a name for this workflow"
	ti.CharLimit = 64
	ti.Width = 40

	ro := textinput.New()
	ro.Placeholder = "Runner Name"
	ro.CharLimit = 64
	ro.Width = 50

	return model{
		state:          stateWorkflowName,
		supportedSched: schedule,
		cronFrequency:  cron,
		supportedCloud: cloud,
		supportedLang:  lang, // Language list component
		textInput:      ti,
		runsOnInput:    ro,
	}
}

// Init initializes the program and starts text input blinking
func (m model) Init() tea.Cmd {
	m.textInput.Focus()
	return textinput.Blink
}
