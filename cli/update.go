package cli

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"workflo/githubactions"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v41/github"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/oauth2"
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

	case stateCustomCron:
		m.textInput, cmd = m.textInput.Update(msg)
		return m.handleCustomCronState(msg, cmd)

	case stateLanguage:
		m.supportedLang, cmd = m.supportedLang.Update(msg)
		return m.handleLanguageState(msg, cmd)

	case stateGitCheckoutOption:
		m.gitCheckoutOption, cmd = m.gitCheckoutOption.Update(msg)
		return m.handleGitCheckoutOptionState(msg, cmd)

	case stateGitBranchSelection:
		m.gitBranchInput.Focus()
		m.gitBranchInput, cmd = m.gitBranchInput.Update(msg)
		return m.handleGitBranchSelectionState(msg, cmd)

	case stateCloudProvider:
		m.supportedCloud, cmd = m.supportedCloud.Update(msg)
		return m.handleCloudProviderState(msg, cmd)

	case stateConfigureAWSCredentials:
		return m.handleConfigureAWSCredentialsState()

	case stateConfigureAWSAccessKeyID:
		m.awsAccessKeyIDInput.Focus()
		m.awsAccessKeyIDInput, cmd = m.awsAccessKeyIDInput.Update(msg)
		return m.handleConfigureAWSAccessKeyIDState(msg, cmd)

	case stateConfigureAWSSecretAccessKey:
		m.awsSecretAccessKeyInput.Focus()
		m.awsSecretAccessKeyInput, cmd = m.awsSecretAccessKeyInput.Update(msg)
		return m.handleConfigureAWSSecretAccessKeyState(msg, cmd)

	case stateConfigureAWSRegion:
		m.awsRegionInput.Focus()
		m.awsRegionInput, cmd = m.awsRegionInput.Update(msg)
		return m.handleConfigureAWSRegionState(msg, cmd)

	case stateConfigureAzureCredentials:
		return m.handleConfigureAzureCredentialsState()

	case stateConfigureAzureClientID:
		m.azureClientIDInput.Focus()
		m.azureClientIDInput, cmd = m.azureClientIDInput.Update(msg)
		return m.handleConfigureAzureClientIDState(msg, cmd)

	case stateConfigureAzureClientSecret:
		m.azureClientSecretInput.Focus()
		m.azureClientSecretInput, cmd = m.azureClientSecretInput.Update(msg)
		return m.handleConfigureAzureClientSecretState(msg, cmd)

	case stateConfigureAzureTenantID:
		m.azureTenantIDInput.Focus()
		m.azureTenantIDInput, cmd = m.azureTenantIDInput.Update(msg)
		return m.handleConfigureAzureTenantIDState(msg, cmd)

	case stateConfigureAzureSubscriptionID:
		m.azureSubscriptionIDInput.Focus()
		m.azureSubscriptionIDInput, cmd = m.azureSubscriptionIDInput.Update(msg)
		return m.handleConfigureAzureSubscriptionIDState(msg, cmd)

	case stateConfigureGCPCredentials:
		return m.handleConfigureGCPCredentialsState()

	case stateConfigureGCPServiceAccountKey:
		m.gcpServiceAccountKeyInput.Focus()
		m.gcpServiceAccountKeyInput, cmd = m.gcpServiceAccountKeyInput.Update(msg)
		return m.handleConfigureGCPServiceAccountKeyState(msg, cmd)

	case stateConfigureGCPProjectID:
		m.gcpProjectIDInput.Focus()
		m.gcpProjectIDInput, cmd = m.gcpProjectIDInput.Update(msg)
		return m.handleConfigureGCPProjectIDState(msg, cmd)

	case stateConfigureSecretsOption:
		m.configureSecretsOption, cmd = m.configureSecretsOption.Update(msg)
		return m.handleConfigureSecretsOptionState(msg, cmd)

	case stateGitHubUsername:
		m.githubUsernameInput.Focus()
		m.githubUsernameInput, cmd = m.githubUsernameInput.Update(msg)
		return m.handleGitHubUsernameState(msg, cmd)

	case stateGitHubRepoName:
		m.githubRepoNameInput.Focus()
		m.githubRepoNameInput, cmd = m.githubRepoNameInput.Update(msg)
		return m.handleGitHubRepoNameState(msg, cmd)

	case stateGitHubToken:
		m.githubTokenInput.Focus()
		m.githubTokenInput, cmd = m.githubTokenInput.Update(msg)
		return m.handleGitHubTokenState(msg, cmd)

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
			stepsYaml := githubactions.GetSkeleton(m.language, m.cloud, m.workflowNameUpper, m.awsRegion)
			steps := githubactions.ParseSteps(stepsYaml)

			// If git checkout is requested, add a step
			if m.gitCheckout {
				checkoutStep := githubactions.Step{
					Name: "Checkout code",
					Uses: "actions/checkout@v2",
					With: map[string]string{
						"ref": m.gitBranch,
					},
				}
				steps = append([]githubactions.Step{checkoutStep}, steps...)
			}

			// Create the job with runner and steps
			job := githubactions.Job{
				RunsOn: m.runsOn,
				Steps:  steps,
			}

			// Add the job to the workflow
			workflow.AddJob("build", job)

			// Generate the YAML file, handling any errors
			// Changed the filename to "workflow.yml"
			err := workflow.GenerateYAML("workflow.yml", true)
			if err != nil {
				fmt.Printf("Error generating workflow YAML: %v\n", err)
			} else {
				fmt.Println("Workflow YAML generated successfully.")
			}

			// Configure secrets if the user chose to
			if m.configureSecrets {
				// Initialize GitHub client
				ctx := context.Background()
				ts := oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: m.githubToken},
				)
				tc := oauth2.NewClient(ctx, ts)
				client := github.NewClient(tc)

				// Configure secrets
				var secrets map[string]string
				if m.cloud == "AWS" {
					secrets = m.awsSecrets
				} else if m.cloud == "Azure" {
					secrets = m.azureSecrets
				} else if m.cloud == "GCP" {
					secrets = m.gcpSecrets
				} else {
					secrets = make(map[string]string)
				}

				err = configureGitHubSecrets(ctx, client, m.githubUsername, m.githubRepoName, secrets)
				if err != nil {
					fmt.Printf("Error configuring GitHub secrets: %v\n", err)
				} else {
					fmt.Println("GitHub secrets configured successfully.")
				}
			}

			// Exit the program
			return m, tea.Quit
		}
	}
	return m, cmd
}

// Function to configure GitHub secrets
func configureGitHubSecrets(ctx context.Context, client *github.Client, owner, repo string, secrets map[string]string) error {
	// Get public key for the repository
	publicKey, _, err := client.Actions.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("error getting public key: %v", err)
	}

	// Encrypt and upload each secret
	for name, value := range secrets {
		// Encrypt the secret value
		encryptedValue, err := encryptSecret([]byte(value), *publicKey.Key)
		if err != nil {
			return fmt.Errorf("error encrypting secret %s: %v", name, err)
		}

		// Create the secret
		secret := &github.EncryptedSecret{
			Name:           name,
			KeyID:          *publicKey.KeyID,
			EncryptedValue: encryptedValue,
		}

		_, err = client.Actions.CreateOrUpdateRepoSecret(ctx, owner, repo, secret)
		if err != nil {
			return fmt.Errorf("error creating/updating secret %s: %v", name, err)
		}
	}

	return nil
}

// Function to encrypt the secret value using the repository's public key
func encryptSecret(secretValue []byte, publicKey string) (string, error) {
	// Decode the public key from Base64
	keyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("error decoding public key: %v", err)
	}

	var publicKeyBytes [32]byte
	copy(publicKeyBytes[:], keyBytes)

	// Encrypt the secret using sealed box
	encryptedBytes, err := box.SealAnonymous(nil, secretValue, &publicKeyBytes, rand.Reader)
	if err != nil {
		return "", fmt.Errorf("error encrypting secret: %v", err)
	}

	// Return the encrypted secret in Base64 encoding
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
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
			m.workflowNameUpper = strings.ToUpper(m.workflowName)
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
					m.customCron = getCronExpression(frequency)
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
			m.textInput.Reset()
			m.state = stateLanguage
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
				m.state = stateGitCheckoutOption
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// handleGitCheckoutOptionState processes input for the Git Checkout Option state
func (m model) handleGitCheckoutOptionState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selectedOption := m.gitCheckoutOption.SelectedItem()
			if selectedOption != nil {
				choice := selectedOption.FilterValue()
				if choice == "Yes" {
					m.gitCheckout = true
					m.state = stateGitBranchSelection
				} else {
					m.gitCheckout = false
					m.state = stateCloudProvider
				}
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// handleGitBranchSelectionState processes input for the Git Branch Selection state
func (m model) handleGitBranchSelectionState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.gitBranch = m.gitBranchInput.Value()
			m.gitBranchInput.Reset()
			m.state = stateCloudProvider
			return m, textinput.Blink
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
				switch m.cloud {
				case "AWS":
					m.state = stateConfigureAWSCredentials
				case "Azure":
					m.state = stateConfigureAzureCredentials
				case "GCP":
					m.state = stateConfigureGCPCredentials
				default:
					m.state = stateConfigureSecretsOption
				}
			} else {
				m.cloud = ""
				m.state = stateConfigureSecretsOption
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// handleConfigureAWSCredentialsState transitions to AWS credential input states
func (m model) handleConfigureAWSCredentialsState() (tea.Model, tea.Cmd) {
	// Transition to the AWS Access Key ID input state
	m.state = stateConfigureAWSAccessKeyID
	return m, textinput.Blink
}

// AWS Credentials States
func (m model) handleConfigureAWSAccessKeyIDState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			secretKey := fmt.Sprintf("%s_AWS_ACCESS_KEY_ID", m.workflowNameUpper)
			m.awsSecrets[secretKey] = m.awsAccessKeyIDInput.Value()
			m.awsAccessKeyIDInput.Reset()
			m.state = stateConfigureAWSSecretAccessKey
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m model) handleConfigureAWSSecretAccessKeyState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			secretKey := fmt.Sprintf("%s_AWS_SECRET_ACCESS_KEY", m.workflowNameUpper)
			m.awsSecrets[secretKey] = m.awsSecretAccessKeyInput.Value()
			m.awsSecretAccessKeyInput.Reset()
			m.state = stateConfigureAWSRegion
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m model) handleConfigureAWSRegionState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			secretKey := fmt.Sprintf("%s_AWS_REGION", m.workflowNameUpper)
			m.awsSecrets[secretKey] = m.awsRegionInput.Value()
			m.awsRegion = m.awsRegionInput.Value()
			m.awsRegionInput.Reset()
			m.state = stateConfigureSecretsOption
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// handleConfigureAzureCredentialsState transitions to Azure credential input states
func (m model) handleConfigureAzureCredentialsState() (tea.Model, tea.Cmd) {
	// Transition to the Azure Client ID input state
	m.state = stateConfigureAzureClientID
	return m, textinput.Blink
}

// Azure Credentials States
func (m model) handleConfigureAzureClientIDState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			secretKey := fmt.Sprintf("%s_AZURE_CLIENT_ID", m.workflowNameUpper)
			m.azureSecrets[secretKey] = m.azureClientIDInput.Value()
			m.azureClientIDInput.Reset()
			m.state = stateConfigureAzureClientSecret
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m model) handleConfigureAzureClientSecretState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			secretKey := fmt.Sprintf("%s_AZURE_CLIENT_SECRET", m.workflowNameUpper)
			m.azureSecrets[secretKey] = m.azureClientSecretInput.Value()
			m.azureClientSecretInput.Reset()
			m.state = stateConfigureAzureTenantID
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m model) handleConfigureAzureTenantIDState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			secretKey := fmt.Sprintf("%s_AZURE_TENANT_ID", m.workflowNameUpper)
			m.azureSecrets[secretKey] = m.azureTenantIDInput.Value()
			m.azureTenantIDInput.Reset()
			m.state = stateConfigureAzureSubscriptionID
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m model) handleConfigureAzureSubscriptionIDState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			secretKey := fmt.Sprintf("%s_AZURE_SUBSCRIPTION_ID", m.workflowNameUpper)
			m.azureSecrets[secretKey] = m.azureSubscriptionIDInput.Value()
			m.azureSubscriptionIDInput.Reset()
			m.state = stateConfigureSecretsOption
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// handleConfigureGCPCredentialsState transitions to GCP credential input states
func (m model) handleConfigureGCPCredentialsState() (tea.Model, tea.Cmd) {
	// Transition to the GCP Service Account Key input state
	m.state = stateConfigureGCPServiceAccountKey
	return m, textinput.Blink
}

// GCP Credentials States
func (m model) handleConfigureGCPServiceAccountKeyState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			secretKey := fmt.Sprintf("%s_GOOGLE_APPLICATION_CREDENTIALS_JSON", m.workflowNameUpper)
			m.gcpSecrets[secretKey] = m.gcpServiceAccountKeyInput.Value()
			m.gcpServiceAccountKeyInput.Reset()
			m.state = stateConfigureGCPProjectID
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m model) handleConfigureGCPProjectIDState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			secretKey := fmt.Sprintf("%s_GCP_PROJECT_ID", m.workflowNameUpper)
			m.gcpSecrets[secretKey] = m.gcpProjectIDInput.Value()
			m.gcpProjectIDInput.Reset()
			m.state = stateConfigureSecretsOption
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// handleConfigureSecretsOptionState processes input for configuring secrets via CLI
func (m model) handleConfigureSecretsOptionState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selectedOption := m.configureSecretsOption.SelectedItem()
			if selectedOption != nil {
				choice := selectedOption.FilterValue()
				if choice == "Yes" {
					m.configureSecrets = true
					m.state = stateGitHubUsername
				} else {
					m.configureSecrets = false
					m.state = stateComplete
				}
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// GitHub Credentials States
func (m model) handleGitHubUsernameState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.githubUsername = m.githubUsernameInput.Value()
			m.githubUsernameInput.Reset()
			m.state = stateGitHubRepoName
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m model) handleGitHubRepoNameState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.githubRepoName = m.githubRepoNameInput.Value()
			m.githubRepoNameInput.Reset()
			m.state = stateGitHubToken
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m model) handleGitHubTokenState(msg tea.Msg, cmd tea.Cmd) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.githubToken = m.githubTokenInput.Value()
			m.githubTokenInput.Reset()
			m.state = stateComplete
			return m, textinput.Blink
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, cmd
}
