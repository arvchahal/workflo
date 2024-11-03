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
	stateSchedule
	stateCronFrequency
	stateCloudProvider
	stateCustomCron
	stateComplete
)

// Model struct to store the state and components
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

// item struct implementing list.Item interface
type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }
