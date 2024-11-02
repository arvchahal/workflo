package githubactions

type Workflow struct {
	Name        string
	Description *string
	Jobs        []Job
	On          map[string]interface{}
}

type Job struct {
	Name   string
	RunsOn string
	Steps  []Step
	Env    map[string]string `yaml:",omitempty"`
}

type Step struct {
	Name string
	Uses string            `yaml:",omitempty"`
	Run  string            `yaml:",omitempty"`
	Env  map[string]string `yaml:",omitempty"`
}

// NewWorkflow initializes a new Workflow
func NewWorkflow(name string) *Workflow {
	return &Workflow{Name: name, On: make(map[string]interface{}), Jobs: []Job{}}
}

// AddJob adds a job to the workflow
func (wf *Workflow) AddJob(job Job) {
	wf.Jobs = append(wf.Jobs, job)
}
