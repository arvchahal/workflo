package githubactions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Checks if a path exists on the filesystem
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// parseSteps converts a YAML string into a slice of Step structs
func parseSteps(stepsYaml string) []Step {
	var steps []Step
	// Split the string by newlines and then parse each individual step
	for _, stepYaml := range strings.Split(stepsYaml, "\n- ") {
		if stepYaml == "" {
			continue
		}

		// Prepend the dash to each step to maintain valid YAML formatting
		stepYaml = "- " + stepYaml
		var step Step
		if err := yaml.Unmarshal([]byte(stepYaml), &step); err == nil {
			steps = append(steps, step)
		}
	}
	return steps
}

/* func buildWf(actionName string, schedule []string, jobs []string  ) *Workflow{
	sdcheudle_type := scheduel[0]
	cron_schedule :=nil
	if len(schedule)>1:
		cron_schedule = schedule[1]
	for job in jobs:


 }
*/
// Generates YAML from a Workflow struct and writes it to a file
func (wf *Workflow) GenerateYAML(filename string, overwrite bool) error {
	dirPath := ".github/workflows"

	// Check if the directory exists, create it if not
	dirExists, err := pathExists(dirPath)
	if err != nil {
		return fmt.Errorf("error checking directory: %v", err)
	}
	if !dirExists {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("error creating github workflows directory: %v", err)
		}
	}
	// Check if file exists in the directory and handle overwrite flag
	filePath := filepath.Join(dirPath, filename)
	if exists, _ := pathExists(filePath); exists && !overwrite {
		return fmt.Errorf("file '%s' already exists and overwrite is set to false; aborting", filename)
	}

	// Marshal the workflow struct into YAML format
	data, err := yaml.Marshal(wf)
	if err != nil {
		return fmt.Errorf("error marshaling YAML: %v", err)
	}

	// Write the YAML data to the file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error writing YAML file: %v", err)
	}

	fmt.Println("Workflow YAML generated successfully.")
	return nil
}
