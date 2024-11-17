package githubactions

import "fmt"

// Basic reusable steps for common GitHub Actions workflows
var BasicSteps = map[string]string{
	//	"checkout": `
	//   - name: Check out the code
	//     uses: actions/checkout@v2`,
}

// Language-specific setup steps
var LanguageSkeletons = map[string]string{
	"Go": `
	- name: Set up Go
  uses: actions/setup-go@v2
  with:
    go-version: '^1.15'
- run: go build -v ./...
- run: go test -v ./...`,

	"Python": `
	- name: Set up Python
 uses: actions/setup-python@v2
  with:
    python-version: '3.x'
- run: pip install -r requirements.txt
- run: pytest`,

	"Node.js": `
	- name: Set up Node.js
 uses: actions/setup-node@v2
  with:
    node-version: '16'
- run: npm install
- run: npm test`,
}

// Cloud provider-specific setup steps
var CloudProviderSkeletons = map[string]string{
	"AWS": `
- name: Configure AWS Credentials
  uses: aws-actions/configure-aws-credentials@v1
  with:
    aws-access-key-id: ${{ secrets.%[1]s_AWS_ACCESS_KEY_ID }}
    aws-secret-access-key: ${{ secrets.%[1]s_AWS_SECRET_ACCESS_KEY }}
    aws-region: ${{ secrets.%[1]s_AWS_REGION }}`,

	"Azure": `
- name: Azure Login
  uses: azure/login@v1
  with:
    creds: ${{ secrets.%[1]s_AZURE_CREDENTIALS }}`,

	"GCP": `
- name: Authenticate to Google Cloud
  uses: google-github-actions/setup-gcloud@v1
  with:
    service_account_key: ${{ secrets.%[1]s_GOOGLE_APPLICATION_CREDENTIALS_JSON }}
    project_id: %[2]s`,
}

// GetSkeleton generates the steps based on selected language and cloud provider
func GetSkeleton(language, cloudProvider, workflowNameUpper string, cloudParam string) string {
	action := ""

	// Add language-specific setup if available
	if langSteps, ok := LanguageSkeletons[language]; ok {
		action += langSteps + "\n"
	}

	// Add cloud provider-specific setup if available
	if cloudSteps, ok := CloudProviderSkeletons[cloudProvider]; ok {
		switch cloudProvider {
		case "AWS", "GCP":
			cloudSteps = fmt.Sprintf(cloudSteps, workflowNameUpper, cloudParam)
		case "Azure":
			cloudSteps = fmt.Sprintf(cloudSteps, workflowNameUpper)
		}
		action += cloudSteps + "\n"
	}

	return action
}
