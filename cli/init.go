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

	// Yes/No options for Git Checkout and Configure Secrets
	yesNoOptions := []list.Item{
		item("Yes"),
		item("No"),
	}

	// Initialize list components
	schedule := list.New(schedulers, list.NewDefaultDelegate(), 50, 15)
	schedule.Title = "Select a scheduling type:"
	schedule.SetShowStatusBar(false)
	schedule.SetShowHelp(false)

	cloud := list.New(cloudProviders, list.NewDefaultDelegate(), 50, 15)
	cloud.Title = "Select a Cloud provider to configure credentials:"
	cloud.SetShowStatusBar(false)
	cloud.SetShowHelp(false)

	cron := list.New(cronOptions, list.NewDefaultDelegate(), 50, 15)
	cron.Title = "Select Cron Frequency:"
	cron.SetShowStatusBar(false)
	cron.SetShowHelp(false)

	gitCheckoutOption := list.New(yesNoOptions, list.NewDefaultDelegate(), 50, 7)
	gitCheckoutOption.Title = "Do you want to perform a git checkout action?"
	gitCheckoutOption.SetShowStatusBar(false)
	gitCheckoutOption.SetShowHelp(false)

	configureSecretsOption := list.New(yesNoOptions, list.NewDefaultDelegate(), 50, 7)
	configureSecretsOption.Title = "Do you want to configure secrets via the CLI?"
	configureSecretsOption.SetShowStatusBar(false)
	configureSecretsOption.SetShowHelp(false)

	// Initialize text inputs
	ti := textinput.New()
	ti.Placeholder = "Enter a name for this workflow"
	ti.CharLimit = 64
	ti.Width = 40

	ro := textinput.New()
	ro.Placeholder = "Runner Name"
	ro.CharLimit = 64
	ro.Width = 50

	// Git branch input
	gb := textinput.New()
	gb.Placeholder = "Enter the branch name to checkout"
	gb.CharLimit = 64
	gb.Width = 40

	// GitHub username input
	githubUsernameInput := textinput.New()
	githubUsernameInput.Placeholder = "Enter your GitHub username"
	githubUsernameInput.CharLimit = 100
	githubUsernameInput.Width = 40

	// GitHub repository name input
	githubRepoNameInput := textinput.New()
	githubRepoNameInput.Placeholder = "Enter your GitHub repository name"
	githubRepoNameInput.CharLimit = 100
	githubRepoNameInput.Width = 40

	// GitHub Personal Access Token input
	githubTokenInput := textinput.New()
	githubTokenInput.Placeholder = "Enter your GitHub Personal Access Token"
	githubTokenInput.CharLimit = 100
	githubTokenInput.Width = 40

	// AWS credential inputs
	awsAccessKeyIDInput := textinput.New()
	awsAccessKeyIDInput.Placeholder = "Enter AWS Access Key ID"
	awsAccessKeyIDInput.CharLimit = 128
	awsAccessKeyIDInput.Width = 50

	awsSecretAccessKeyInput := textinput.New()
	awsSecretAccessKeyInput.Placeholder = "Enter AWS Secret Access Key"
	awsSecretAccessKeyInput.CharLimit = 128
	awsSecretAccessKeyInput.Width = 50

	awsRegionInput := textinput.New()
	awsRegionInput.Placeholder = "Enter AWS Region"
	awsRegionInput.CharLimit = 64
	awsRegionInput.Width = 50

	// Azure credential inputs
	azureClientIDInput := textinput.New()
	azureClientIDInput.Placeholder = "Enter Azure Client ID"
	azureClientIDInput.CharLimit = 128
	azureClientIDInput.Width = 50

	azureClientSecretInput := textinput.New()
	azureClientSecretInput.Placeholder = "Enter Azure Client Secret"
	azureClientSecretInput.CharLimit = 128
	azureClientSecretInput.Width = 50

	azureTenantIDInput := textinput.New()
	azureTenantIDInput.Placeholder = "Enter Azure Tenant ID"
	azureTenantIDInput.CharLimit = 128
	azureTenantIDInput.Width = 50

	azureSubscriptionIDInput := textinput.New()
	azureSubscriptionIDInput.Placeholder = "Enter Azure Subscription ID"
	azureSubscriptionIDInput.CharLimit = 128
	azureSubscriptionIDInput.Width = 50

	// GCP credential inputs
	gcpServiceAccountKeyInput := textinput.New()
	gcpServiceAccountKeyInput.Placeholder = "Enter GCP Service Account Key (JSON)"
	gcpServiceAccountKeyInput.CharLimit = 5000
	gcpServiceAccountKeyInput.Width = 50

	gcpProjectIDInput := textinput.New()
	gcpProjectIDInput.Placeholder = "Enter GCP Project ID"
	gcpProjectIDInput.CharLimit = 128
	gcpProjectIDInput.Width = 50

	return model{
		state:                     stateWorkflowName,
		supportedSched:            schedule,
		cronFrequency:             cron,
		supportedCloud:            cloud,
		supportedLang:             lang,
		textInput:                 ti,
		runsOnInput:               ro,
		gitCheckoutOption:         gitCheckoutOption,
		configureSecretsOption:    configureSecretsOption,
		gitBranchInput:            gb,
		githubUsernameInput:       githubUsernameInput,
		githubRepoNameInput:       githubRepoNameInput,
		githubTokenInput:          githubTokenInput,
		awsAccessKeyIDInput:       awsAccessKeyIDInput,
		awsSecretAccessKeyInput:   awsSecretAccessKeyInput,
		awsRegionInput:            awsRegionInput,
		azureClientIDInput:        azureClientIDInput,
		azureClientSecretInput:    azureClientSecretInput,
		azureTenantIDInput:        azureTenantIDInput,
		azureSubscriptionIDInput:  azureSubscriptionIDInput,
		gcpServiceAccountKeyInput: gcpServiceAccountKeyInput,
		gcpProjectIDInput:         gcpProjectIDInput,
		awsSecrets:                make(map[string]string),
		azureSecrets:              make(map[string]string),
		gcpSecrets:                make(map[string]string),
	}
}

// Init initializes the program and starts text input blinking
func (m model) Init() tea.Cmd {
	m.textInput.Focus()
	return textinput.Blink
}
