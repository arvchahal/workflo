package githubactions

type workflow struct {
	name        string
	description *string
	jobs        []job
	on          map[string]interface{}
}

type job struct {
	name  string
	on    string
	steps []step
	Env   map[string]string `yaml:",omitempty"`
}
type step struct {
	Name string
	Uses string            `yaml:",omitempty"`
	Run  string            `yaml:",omitempty"`
	Env  map[string]string `yaml:",omitempty"`
}

func newWorkflow(name string) *workflow {
	return &workflow{name: name, on: make(map[string]interface{}), jobs: []job{}}
}

func (wf *workflow) add_job(job job) {
	wf.jobs = append(wf.jobs, job)
}
