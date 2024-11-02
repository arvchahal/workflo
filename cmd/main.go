package main

import (
	"log"

	"github.com/arvchahal/workflo/githubactions"
)

func main() {
	wf := githubactions.NewWorkflow("CI Workflow")
	wf.On["push"] = nil
	wf.On["pull_request"] = nil

	job := githubactions.Job{
		Name:   "Build and Test",
		RunsOn: "ubuntu-latest",
		Steps: []githubactions.Step{
			githubactions.Step{Name: "Checkout", Uses: "actions/checkout@v2"},
			githubactions.Step{Name: "Setup Go", Uses: "actions/setup-go@v2", Env: map[string]string{"go-version": "1.17"}},
			githubactions.Step{Name: "Run Tests", Run: "go test ./..."},
		},
	}
	wf.AddJob(job)

	if err := wf.GenerateYAML("ci.yml", true); err != nil {
		log.Fatalf("Failed to generate YAML: %v", err)
	}

	log.Println("Workflow YAML generated successfully.")
}
