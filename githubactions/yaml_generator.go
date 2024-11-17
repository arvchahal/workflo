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
func ParseSteps(stepsYaml string) []Step {
	var steps []Step
	// Remove any leading/trailing whitespace
	stepsYaml = strings.TrimSpace(stepsYaml)

	// Ensure stepsYaml starts with '- '
	if !strings.HasPrefix(stepsYaml, "- ") {
		stepsYaml = "- " + stepsYaml
	}

	// Parse the stepsYaml into the steps slice
	if err := yaml.Unmarshal([]byte(stepsYaml), &steps); err != nil {
		fmt.Printf("Error parsing steps: %v\n", err)
		return nil
	}

	return steps
}

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

	// Construct the file path
	filePath := filepath.Join(dirPath, filename)

	// Check if file exists in the directory and handle overwrite flag
	if exists, _ := pathExists(filePath); exists && !overwrite {
		return fmt.Errorf("file '%s' already exists and overwrite is set to false; aborting", filename)
	}

	// Marshal the workflow struct into YAML format
	data, err := yaml.Marshal(wf)
	if err != nil {
		return fmt.Errorf("error marshaling YAML: %v", err)
	}

	// Replace `"on":` with `on:` (YAML syntax)
	yamlString := strings.Replace(string(data), `"on":`, "on:", 1)

	// Write the YAML data to the file
	if err := os.WriteFile(filePath, []byte(yamlString), 0644); err != nil {
		return fmt.Errorf("error writing YAML file: %v", err)
	}

	fmt.Println("Workflow YAML generated successfully.")
	return nil
}
