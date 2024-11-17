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

	case stateCustomCron:
		return fmt.Sprintf("Enter custom cron schedule:\n\n%s\n\n(Press Enter to continue)", m.textInput.View())

	case stateLanguage:
		return m.supportedLang.View()

	case stateGitCheckoutOption:
		return m.gitCheckoutOption.View()

	case stateGitBranchSelection:
		return fmt.Sprintf("Enter the branch name to checkout:\n\n%s\n\n(Press Enter to continue)", m.gitBranchInput.View())

	case stateCloudProvider:
		return m.supportedCloud.View()

	case stateConfigureSecretsOption:
		return m.configureSecretsOption.View()

	// Add cases for the intermediate cloud configuration states
	case stateConfigureAWSCredentials:
		return "Preparing to configure AWS credentials..."

	case stateConfigureAzureCredentials:
		return "Preparing to configure Azure credentials..."

	case stateConfigureGCPCredentials:
		return "Preparing to configure GCP credentials..."

	case stateConfigureAWSAccessKeyID:
		return fmt.Sprintf("Enter AWS Access Key ID:\n\n%s\n\n(Press Enter to continue)", m.awsAccessKeyIDInput.View())

	case stateConfigureAWSSecretAccessKey:
		return fmt.Sprintf("Enter AWS Secret Access Key:\n\n%s\n\n(Press Enter to continue)", m.awsSecretAccessKeyInput.View())

	case stateConfigureAWSRegion:
		return fmt.Sprintf("Enter AWS Region:\n\n%s\n\n(Press Enter to continue)", m.awsRegionInput.View())

	case stateConfigureAzureClientID:
		return fmt.Sprintf("Enter Azure Client ID:\n\n%s\n\n(Press Enter to continue)", m.azureClientIDInput.View())

	case stateConfigureAzureClientSecret:
		return fmt.Sprintf("Enter Azure Client Secret:\n\n%s\n\n(Press Enter to continue)", m.azureClientSecretInput.View())

	case stateConfigureAzureTenantID:
		return fmt.Sprintf("Enter Azure Tenant ID:\n\n%s\n\n(Press Enter to continue)", m.azureTenantIDInput.View())

	case stateConfigureAzureSubscriptionID:
		return fmt.Sprintf("Enter Azure Subscription ID:\n\n%s\n\n(Press Enter to continue)", m.azureSubscriptionIDInput.View())

	case stateConfigureGCPServiceAccountKey:
		return fmt.Sprintf("Enter GCP Service Account Key (JSON):\n\n%s\n\n(Press Enter to continue)", m.gcpServiceAccountKeyInput.View())

	case stateConfigureGCPProjectID:
		return fmt.Sprintf("Enter GCP Project ID:\n\n%s\n\n(Press Enter to continue)", m.gcpProjectIDInput.View())

	case stateGitHubUsername:
		return fmt.Sprintf("Enter your GitHub username:\n\n%s\n\n(Press Enter to continue)", m.githubUsernameInput.View())

	case stateGitHubRepoName:
		return fmt.Sprintf("Enter your GitHub repository name:\n\n%s\n\n(Press Enter to continue)", m.githubRepoNameInput.View())

	case stateGitHubToken:
		return fmt.Sprintf("Enter your GitHub Personal Access Token (with 'repo' scope):\n\n%s\n\n(Press Enter to continue)", m.githubTokenInput.View())

	case stateComplete:
		return "Workflow setup completed! Press Enter or Ctrl+C to exit."

	default:
		return "An unexpected error occurred."
	}
}
