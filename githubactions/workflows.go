// githubactions/workflow.go

package githubactions

type Workflow struct {
	Name        string                 `yaml:"name"`
	Description *string                `yaml:"description,omitempty"`
	On          map[string]interface{} `yaml:"on"`
	Jobs        map[string]Job         `yaml:"jobs"`
}

type Job struct {
	RunsOn string            `yaml:"runs-on"`
	Steps  []Step            `yaml:"steps"`
	Env    map[string]string `yaml:"env,omitempty"`
}

type Step struct {
	Name string            `yaml:"name,omitempty"`
	Uses string            `yaml:"uses,omitempty"`
	Run  string            `yaml:"run,omitempty"`
	Env  map[string]string `yaml:"env,omitempty"`
	With map[string]string `yaml:"with,omitempty"`
}

// NewWorkflow initializes a new Workflow
func NewWorkflow(name string) *Workflow {
	return &Workflow{
		Name: name,
		On:   make(map[string]interface{}),
		Jobs: make(map[string]Job),
	}
}

// AddJob adds a job to the workflow
func (wf *Workflow) AddJob(jobName string, job Job) {
	wf.Jobs[jobName] = job
}
