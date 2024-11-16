// update.go
package cli

import (
	"fmt"
	"workflo/githubactions"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateWorkflowName:
		m.textInput.Focus()
		m.textInput, cmd = m.textInput.Update(msg)
		return m.handleWorkflowNameState(msg, cmd)

	case stateRunner:
		m.runsOnInput.Focus()
		m.runsOnInput, cmd = m.runsOnInput.Update(msg)
		return m.handleRunnerState(msg, cmd)

	case stateSchedule:
		m.supportedSched, cmd = m.supportedSched.Update(msg)
		return m.handleScheduleState(msg, cmd)

	case stateCronFrequency:
		m.cronFrequency, cmd = m.cronFrequency.Update(msg)
		return m.handleCronFrequencyState(msg, cmd)

	case stateLanguage:
		m.supportedLang, cmd = m.supportedLang.Update(msg)
		return m.handleLanguageState(msg, cmd)

	case stateCloudProvider:
		m.supportedCloud, cmd = m.supportedCloud.Update(msg)
		return m.handleCloudProviderState(msg, cmd)

	case stateCustomCron:
		m.textInput, cmd = m.textInput.Update(msg)
		return m.handleCustomCronState(msg, cmd)

	case stateComplete:
		return m.handleCompleteState(msg, cmd)
	}
	return m, cmd
}

// handleCompleteState handles the final state where the program exits and generates the YAML file.
func (m model) handleCompleteState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "ctrl+c", "q":
			// Initialize the workflow with the collected inputs
			workflow := githubactions.NewWorkflow(m.workflowName)

			// Map user-friendly names to valid GitHub Actions event names
			eventName := ""
			switch m.schedule {
			case "On dispatch":
				eventName = "workflow_dispatch"
				workflow.On[eventName] = map[string]interface{}{}
			case "On Pull":
				eventName = "pull_request"
				workflow.On[eventName] = map[string]interface{}{
					"branches": []string{"main"},
				}
			case "On Push":
				eventName = "push"
				workflow.On[eventName] = map[string]interface{}{
					"branches": []string{"main"},
				}
			case "Cron Schedule":
				eventName = "schedule"
				cronExpression := m.customCron
				if cronExpression == "" {
					cronExpression = getCronExpression(m.cronFrequency.SelectedItem().FilterValue())
				}
				workflow.On[eventName] = []map[string]string{
					{"cron": cronExpression},
				}
			default:
				fmt.Println("Invalid schedule type selected.")
				return m, tea.Quit
			}

			// Check for empty `runsOn` and default to `ubuntu-latest`
			if m.runsOn == "" {
				m.runsOn = "ubuntu-latest"
			}

			// Generate steps for the job based on language and cloud provider
			stepsYaml := githubactions.GetSkeleton(m.language, m.cloud, true)
			steps := githubactions.ParseSteps(stepsYaml)

			// Create the job with runner and steps, then add to workflow
			job := githubactions.Job{
				RunsOn: m.runsOn,
				Steps:  steps,
			}
			workflow.AddJob("build", job)

			// Generate the YAML file, handling any errors
			err := workflow.GenerateYAML("workflow.yml", true)
			if err != nil {
				fmt.Printf("Error generating workflow YAML: %v\n", err)
			} else {
				fmt.Println("Workflow YAML generated successfully.")
			}

			// Exit the program
			return m, tea.Quit
		}
	}
	return m, cmd
}

// Helper function to map cron frequency to cron expressions
func getCronExpression(frequency string) string {
	switch frequency {
	case "Once a day":
		return "0 0 * * *"
	case "Once a week":
		return "0 0 * * 0"
	case "Once a month":
		return "0 0 1 * *"
	case "Once a year":
		return "0 0 1 1 *"
	default:
		return "0 0 * * *" // Default to once a day
	}
}

// handleWorkflowNameState processes input for the Workflow Name state
func (m model) handleWorkflowNameState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.workflowName = m.textInput.Value()
			m.textInput.Reset()
			m.state = stateRunner
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// handleRunnerState processes input for the Runner state
func (m model) handleRunnerState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.runsOn = m.runsOnInput.Value()
			m.runsOnInput.Reset()
			m.state = stateSchedule
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// handleScheduleState processes input for the Schedule state
func (m model) handleScheduleState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
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
					m.state = stateLanguage
				}
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// handleCronFrequencyState processes input for the Cron Frequency state
func (m model) handleCronFrequencyState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selectedFrequency := m.cronFrequency.SelectedItem()
			if selectedFrequency != nil {
				frequency := selectedFrequency.FilterValue()
				if frequency == "Other (Enter custom cron)" {
					m.textInput.Placeholder = "Enter custom cron schedule"
					m.textInput.SetValue("")
					m.textInput.Focus()
					m.state = stateCustomCron
					return m, textinput.Blink
				} else {
					m.schedule = frequency
					m.state = stateLanguage
				}
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// handleCustomCronState processes input for the Custom Cron state
func (m model) handleCustomCronState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Set the custom cron value from user input
			m.customCron = m.textInput.Value()
			m.schedule = m.customCron // Use custom cron as the schedule
			m.state = stateLanguage   // Transition to stateLanguage
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// handleLanguageState processes input for the Language state
func (m model) handleLanguageState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selectedLang := m.supportedLang.SelectedItem()
			if selectedLang != nil {
				m.language = selectedLang.FilterValue()
				m.state = stateCloudProvider
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// handleCloudProviderState processes input for the Cloud Provider state
func (m model) handleCloudProviderState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selectedCloud := m.supportedCloud.SelectedItem()
			if selectedCloud != nil {
				m.cloud = selectedCloud.FilterValue()
				m.state = stateComplete
			} else {
				m.cloud = ""
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}
