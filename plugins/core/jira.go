package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
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
	case "assigned-issues":
		if len(args) < 2 {
			return fmt.Errorf("usage: assigned-issues <userEmail>")
		}
		return p.assignedIssues(args[1])
	case "summary":
		if len(args) < 2 {
			return fmt.Errorf("usage: summary <userEmail>")
		}

		// Fetch assigned issues
		issues, err := p.fetchAssignedIssues(args[1])
		if err != nil {
			return err
		}

		// Call AI Summary function
		return p.generateAISummary(args[1], issues)
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

func (p *JiraPlugin) assignedIssues(userEmail string) error {
	jql := fmt.Sprintf("assignee='%s' ORDER BY updated DESC", userEmail)
	encodedJQL := url.QueryEscape(jql)

	url := fmt.Sprintf("%s/rest/api/3/search?jql=%s", p.baseURL, encodedJQL)

	fmt.Println("Final Request URL:", url) // Debugging output

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.SetBasicAuth(p.email, p.apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second} // Increased timeout
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to fetch assigned issues, status: %s, response: %s", resp.Status, string(body))
	}

	// Parse response
	var result struct {
		Issues []struct {
			Key    string `json:"key"`
			Fields struct {
				Summary   string `json:"summary"`
				IssueType struct {
					Name string `json:"name"`
				} `json:"issuetype"`
				Status struct {
					Name string `json:"name"`
				} `json:"status"`
				Updated string `json:"updated"`
			} `json:"fields"`
		} `json:"issues"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if len(result.Issues) == 0 {
		fmt.Printf("📌 No issues assigned to %s.\n", userEmail)
		return nil
	}

	fmt.Printf("📌 Issues assigned to %s:\n", userEmail)
	for _, issue := range result.Issues {
		fmt.Printf("- %s [%s]: %s (Status: %s, Updated: %s)\n",
			issue.Key, issue.Fields.IssueType.Name, issue.Fields.Summary, issue.Fields.Status.Name, issue.Fields.Updated)
	}

	return nil
}

func (p *JiraPlugin) generateAISummary(userEmail string, issues []struct {
	Key       string
	Summary   string
	IssueType string
	Status    string
	Updated   string
}) error {

	apiKey := os.Getenv("OPENAI_API_KEY")
	orgID := os.Getenv("OPENAI_ORG")

	if apiKey == "" {
		return fmt.Errorf("OPENAI_API_KEY is not set. Please export your API key.")
	}
	formattedIssuesText := ""
	for _, issue := range issues {
		formattedIssuesText += fmt.Sprintf("- [%s] %s: %s (Status: %s, Updated: %s)\n",
			issue.IssueType, issue.Key, issue.Summary, issue.Status, issue.Updated)
	}

	prompt := fmt.Sprintf(`
      I am preparing a self-evaluation for my work. Please summarize my Jira contributions in a professional yet concise way.
      Focus on the impact of my work rather than just listing tasks. Frame the summary as if I am describing my achievements for a performance review.

      Here are my recent Jira contributions:
      %s`, formattedIssuesText)

	payload := map[string]interface{}{
		//"model": "gpt-4",
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "You are an assistant summarizing Jira issue contributions."},
			{"role": "user", "content": prompt},
		},
		"max_tokens": 200,
	}

	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	//orgid header
	if orgID != "" {
		req.Header.Set("OpenAI-Organization", orgID)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get AI summary, status: %s, response: %s", resp.Status, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse AI response: %v", err)
	}

	if len(result.Choices) > 0 {
		fmt.Println("\n📌 AI-Generated Summary:\n" + result.Choices[0].Message.Content)
	} else {
		fmt.Println("\n⚠️ AI did not return a summary.")
	}

	return nil
}

func (p *JiraPlugin) fetchAssignedIssues(userEmail string) ([]struct {
	Key       string
	Summary   string
	IssueType string
	Status    string
	Updated   string
}, error) {
	// Properly encode JQL query
	jql := fmt.Sprintf("assignee='%s' ORDER BY updated DESC", userEmail)
	encodedJQL := url.QueryEscape(jql)

	// Construct request URL
	url := fmt.Sprintf("%s/rest/api/3/search?jql=%s", p.baseURL, encodedJQL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.SetBasicAuth(p.email, p.apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch assigned issues, status: %s, response: %s", resp.Status, string(body))
	}

	// Parse response
	var result struct {
		Issues []struct {
			Key    string `json:"key"`
			Fields struct {
				Summary   string `json:"summary"`
				IssueType struct {
					Name string `json:"name"`
				} `json:"issuetype"`
				Status struct {
					Name string `json:"name"`
				} `json:"status"`
				Updated string `json:"updated"`
			} `json:"fields"`
		} `json:"issues"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Convert to expected return format
	var formattedIssues []struct {
		Key       string
		Summary   string
		IssueType string
		Status    string
		Updated   string
	}
	for _, issue := range result.Issues {
		formattedIssues = append(formattedIssues, struct {
			Key       string
			Summary   string
			IssueType string
			Status    string
			Updated   string
		}{
			Key:       issue.Key,
			Summary:   issue.Fields.Summary,
			IssueType: issue.Fields.IssueType.Name,
			Status:    issue.Fields.Status.Name,
			Updated:   issue.Fields.Updated,
		})
	}

	return formattedIssues, nil
}
