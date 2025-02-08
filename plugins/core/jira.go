package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// JiraPlugin represents the Jira integration plugin
type JiraPlugin struct {
	baseURL  string
	apiToken string
	email    string
}

func (p *JiraPlugin) Init() error {
	fmt.Println("Initializing Jira plugin...")
	p.baseURL = os.Getenv("JIRA_BASE_URL")
	p.email = os.Getenv("JIRA_EMAIL")
	p.apiToken = os.Getenv("JIRA_API_TOKEN")

	if p.baseURL == "" || p.apiToken == "" || p.email == "" {
		return fmt.Errorf("JIRA_BASE_URL, JIRA_API_TOKEN, or JIRA_EMAIL is not set in environment variables")
	}

	fmt.Println("Jira plugin initialized with base URL:", p.baseURL)
	fmt.Printf("Using email: %s, API token length: %d\n", p.email, len(p.apiToken))
	return nil
}

func (p *JiraPlugin) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided")
	}

	switch args[0] {
	case "create-issue":
		if len(args) < 3 {
			return fmt.Errorf("usage: create-issue <summary> <description>")
		}
		return p.createIssue(args[1], args[2])
	case "list-issues":
		if len(args) < 2 {
			return fmt.Errorf("usage: list-issues <projectKey>")
		}
		return p.listIssues(args[1])
	default:
		return fmt.Errorf("unknown Jira command: %s", args[0])
	}
}

func (p *JiraPlugin) Info() (string, string) {
	return "jira", "Integration with Jira for tracking issues"
}

func (p *JiraPlugin) createIssue(summary, description string) error {
	url := p.baseURL + "/rest/api/3/issue"
	payload := map[string]interface{}{
		"fields": map[string]interface{}{
			"project": map[string]string{"key": "SCRUM"}, // Replace with your project key
			"summary": summary,
			"description": map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []map[string]interface{}{
					{
						"type": "paragraph",
						"content": []map[string]interface{}{
							{
								"type": "text",
								"text": description,
							},
						},
					},
				},
			},
			"issuetype": map[string]string{"name": "Task"},
		},
	}

	jsonPayload, _ := json.Marshal(payload)

	fmt.Printf("Payload: %s\n", string(jsonPayload)) // Print the payload
	fmt.Printf("Request URL: %s\n", url)             // Print the URL

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.SetBasicAuth(p.email, p.apiToken)

	req.Header.Set("Content-Type", "application/json")

	fmt.Println("Request Headers:", req.Header)

	fmt.Println("Authorization Header:", req.Header.Get("Authorization"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error response: %s\n", string(body)) // Log the full response body
		return fmt.Errorf("failed to create issue, status: %s, response: %s", resp.Status, string(body))
	}

	fmt.Println("Jira issue created successfully!")
	return nil
}

func (p *JiraPlugin) listIssues(projectKey string) error {
	url := fmt.Sprintf("%s/rest/api/3/search?jql=project=%s", p.baseURL, projectKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.SetBasicAuth(p.email, p.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to fetch issues, status: %s, response: %s", resp.Status, string(body))
	}

	// Parse the response body
	var result struct {
		Issues []struct {
			Key    string `json:"key"`
			Fields struct {
				Summary string `json:"summary"`
			} `json:"fields"`
		} `json:"issues"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	// Display the issues
	fmt.Printf("Issues for project %s:\n", projectKey)
	for _, issue := range result.Issues {
		fmt.Printf("- %s: %s\n", issue.Key, issue.Fields.Summary)
	}

	return nil
}
