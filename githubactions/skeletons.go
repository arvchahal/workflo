package githubactions

// import "fmt"

// Basic reusable steps for common GitHub Actions workflows
var BasicSteps = map[string]string{
	"checkout": `
- name: Check out the code
  uses: actions/checkout@v2`,

	"setup-node": `
- name: Set up Node.js
  uses: actions/setup-node@v2
  with:
    node-version: '16'`,

	"setup-python": `
- name: Set up Python
  uses: actions/setup-python@v2
  with:
    python-version: '3.x'`,

	"setup-go": `
- name: Set up Go
  uses: actions/setup-go@v2
  with:
    go-version: '^1.15'`,
}

// Language-specific setup steps
var LanguageSkeletons = map[string]string{
	"Go": `
- run: go build -v ./...
- run: go test -v ./...`,

	"Python": `
- run: pip install -r requirements.txt
- run: pytest`,

	"Node.js": `
- run: npm install
- run: npm test`,
}

// Cloud provider-specific setup steps, including SSO configurations
var CloudProviderSkeletons = map[string]string{
	"AWS": `
- name: Configure AWS SSO credentials
  uses: aws-actions/configure-aws-credentials@v1
  with:
    role-to-assume: arn:aws:iam::123456789012:role/SSOReadOnly
    aws-region: us-west-2`,

	"Azure": `
- name: Azure Login
  uses: azure/login@v1
  with:
    creds: ${{ secrets.AZURE_CREDENTIALS }}`,

	"GCP": `
- name: Authenticate to Google Cloud
  uses: google-github-actions/auth@v1
  with:
    credentials_json: ${{ secrets.GCP_CREDENTIALS }}`,

	"GCP-Workload": `
- name: Authenticate to Google Cloud with Workload Identity Federation
  uses: google-github-actions/auth@v1
  with:
    workload_identity_provider: "projects/123456789/locations/global/workloadIdentityPools/my-pool/providers/my-provider"
    service_account: "my-service-account@my-project.iam.gserviceaccount.com"`,
}

// GetSkeleton generates the GitHub Actions workflow based on selected language, cloud provider, and optional setup steps
func GetSkeleton(language, cloudProvider string, includeSSO bool) string {
	// Start with basic steps like checkout
	action := BasicSteps["checkout"] + "\n"

	// Add language-specific setup if available
	if langSteps, ok := LanguageSkeletons[language]; ok {
		action += langSteps + "\n"
	}

	// Add cloud provider-specific setup, including SSO where applicable
	if cloudSteps, ok := CloudProviderSkeletons[cloudProvider]; ok {
		action += cloudSteps + "\n"
	}

	// Optionally, include additional setup steps based on selected options
	if includeSSO && cloudProvider == "GCP" {
		action += CloudProviderSkeletons["GCP-Workload"] + "\n"
	} else if includeSSO && cloudProvider == "AWS" {
		action += CloudProviderSkeletons["AWS"] + "\n"
	}

	return action
}
