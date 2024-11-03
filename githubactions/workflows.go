package githubactions

type Workflow struct {
	Name        string
	Description *string
	On          map[string]interface{}
	Jobs        map[string]Job
}

type Job struct {
	RunsOn string
	Steps  []Step
	Env    map[string]string `yaml:",omitempty"`
}

type Step struct {
	Name string            `yaml:",omitempty"`
	Uses string            `yaml:",omitempty"`
	Run  string            `yaml:",omitempty"`
	Env  map[string]string `yaml:",omitempty"`
	With map[string]string `yaml:",omitempty"`
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
