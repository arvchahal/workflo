// view.go
package cli

import "fmt"

// View renders the UI based on the current state
func (m model) View() string {
	switch m.state {
	case stateWorkflowName:
		return fmt.Sprintf("Enter a name for this workflow:\n\n%s\n\n(Press Enter to continue)", m.textInput.View())
	case stateRunner:
		return fmt.Sprintf("Enter a specific Runner for the workflow file (leave empty for default: ubuntu-latest): \n\n%s\n\n(Press Enter to continue)", m.runsOnInput.View())
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
