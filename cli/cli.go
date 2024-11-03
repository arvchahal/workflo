package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"workflo/githubactions"
)

var (
	titleStyle         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5733"))
	selectedTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#33FF57"))
	descriptionStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
)

func genCustomList() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = selectedTitleStyle
	delegate.Styles.NormalTitle = titleStyle
	delegate.Styles.NormalDesc = descriptionStyle
	return delegate
}

type state int

const (
	stateProjectType state = iota
	stateWorkflowName
	stateFinalize
)

type model struct {
	state       state
	list        list.Model
	textInput   textinput.Model
	workflow    *githubactions.Workflow
	projectType string
	cursor      int
}

func NewModel() model {
	// Define choices for project type
	items := []list.Item{
		item("Go"),
		item("Python"),
		item("Node.js"),
		item("None of the Above (default Yaml config)"),
	}

	// Initialize the list component
	//delegate := genCustomList()
	l := list.New(items, list.NewDefaultDelegate(), 100, 15)
	l.Title = "What programming language is your project using?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.DisableQuitKeybindings()
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	ti := textinput.New()
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("")).Background(lipgloss.Color("")) // Invisible cursor
	ti.Blur()                                                                                           // Unfocus to hide cursor initially
	return model{
		state:    stateProjectType,
		list:     l,
		workflow: githubactions.NewWorkflow(""),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case stateProjectType:
		m.list, cmd = m.list.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				selectedItem := m.list.SelectedItem()
				if selectedItem != nil {
					m.projectType = selectedItem.FilterValue()
					m.state = stateWorkflowName

					// Initialize text input for workflow name
					m.textInput = textinput.New()
					m.textInput.Placeholder = "Enter a name for this workflow"
					m.textInput.Focus()
					m.textInput.CharLimit = 64
					m.textInput.Width = 40
				}
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}

		return m, cmd

	case stateWorkflowName:
		m.textInput, cmd = m.textInput.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.workflow.Name = m.textInput.Value()
				m.state = stateFinalize
				// Proceed to finalize or continue to next steps
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}

		return m, cmd

	case stateFinalize:
		// Generate the YAML file
		err := m.workflow.GenerateYAML("ci-workflow.yml", false)
		if err != nil {
			return m, tea.Quit
		}
		fmt.Println("Workflow YAML generated successfully!")
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	switch m.state {
	case stateProjectType:
		return m.list.View()
	case stateWorkflowName:
		return fmt.Sprintf(
			"Enter a name for this workflow:\n\n%s\n\n(Press Enter to continue)",
			m.textInput.View(),
		)
	case stateFinalize:
		return "Workflow generated successfully!\n"
	default:
		return "An unexpected error occurred.\n"
	}
}

// item implements the list.Item interface
type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }
