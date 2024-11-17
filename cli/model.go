// model.go
package cli

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

// Define the application states
type state int

const (
	stateWorkflowName state = iota
	stateRunner
	stateSchedule
	stateCronFrequency
	stateCustomCron
	stateLanguage
	stateGitCheckoutOption
	stateGitBranchSelection
	stateCloudProvider
	stateConfigureAWSCredentials
	stateConfigureAWSAccessKeyID
	stateConfigureAWSSecretAccessKey
	stateConfigureAWSRegion
	stateConfigureAzureCredentials
	stateConfigureAzureClientID
	stateConfigureAzureClientSecret
	stateConfigureAzureTenantID
	stateConfigureAzureSubscriptionID
	stateConfigureGCPCredentials
	stateConfigureGCPServiceAccountKey
	stateConfigureGCPProjectID
	stateConfigureSecretsOption
	stateGitHubUsername
	stateGitHubRepoName
	stateGitHubToken
	stateComplete
)

// Model struct to store the state and components
type model struct {
	state                     state
	supportedSched            list.Model
	supportedCloud            list.Model
	cronFrequency             list.Model
	supportedLang             list.Model
	gitCheckoutOption         list.Model
	configureSecretsOption    list.Model
	textInput                 textinput.Model
	runsOnInput               textinput.Model
	gitBranchInput            textinput.Model
	githubUsernameInput       textinput.Model
	githubRepoNameInput       textinput.Model
	githubTokenInput          textinput.Model
	awsAccessKeyIDInput       textinput.Model
	awsSecretAccessKeyInput   textinput.Model
	awsRegionInput            textinput.Model
	azureClientIDInput        textinput.Model
	azureClientSecretInput    textinput.Model
	azureTenantIDInput        textinput.Model
	azureSubscriptionIDInput  textinput.Model
	gcpServiceAccountKeyInput textinput.Model
	gcpProjectIDInput         textinput.Model
	workflowName              string
	workflowNameUpper         string
	schedule                  string
	cloud                     string
	language                  string
	awsRegion                 string
	customCron                string
	runsOn                    string
	gitCheckout               bool
	gitBranch                 string
	awsSecrets                map[string]string
	azureSecrets              map[string]string
	gcpSecrets                map[string]string
	configureSecrets          bool
	githubUsername            string
	githubRepoName            string
	githubToken               string
}

// item struct implementing list.Item interface
type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }
