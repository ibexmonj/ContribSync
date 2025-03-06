package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ibexmonj/ContribSync/pkg/logger"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// JiraPlugin represents the Jira integration plugin
type JiraPlugin struct {
	baseURL  string
	apiToken string
	email    string
}

func (p *JiraPlugin) LoadEnvVars() error {
	p.baseURL = os.Getenv("JIRA_BASE_URL")
	p.email = os.Getenv("JIRA_EMAIL")
	p.apiToken = os.Getenv("JIRA_API_TOKEN")

	if p.baseURL == "" || p.apiToken == "" || p.email == "" {
		return fmt.Errorf("missing required environment variables")
	}
	return nil
}
func toJSON(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("JSON marshaling error: %w", err)
	}
	return data, nil
}

func wrapError(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

func (p *JiraPlugin) Init() error {
	logger.Logger.Info().Msg("Initializing Jira plugin...")
	if err := p.LoadEnvVars(); err != nil {
		return err
	}
	logger.Logger.Info().Str("Base URL", p.baseURL).Msg("Jira plugin initialized")
	return nil
}

func (p *JiraPlugin) makeRequest(method, endpoint string, body io.Reader, isOpenAI bool) (*http.Response, error) {
	var fullURL string
	reqHeaders := make(map[string]string)

	if strings.HasPrefix(endpoint, "http") {
		fullURL = endpoint
	} else {
		fullURL = p.baseURL + endpoint
	}

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if isOpenAI {
		apiKey := os.Getenv("OPENAI_API_KEY")
		orgID := os.Getenv("OPENAI_ORG")

		if apiKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY is not set. Please export your API key.")
		}

		reqHeaders["Authorization"] = "Bearer " + apiKey
		reqHeaders["Content-Type"] = "application/json"
		if orgID != "" {
			reqHeaders["OpenAI-Organization"] = orgID
		}
	} else {
		reqHeaders["Authorization"] = "Basic " + p.apiToken
		reqHeaders["Content-Type"] = "application/json"
	}

	for key, value := range reqHeaders {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}

func HandleResponseBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		fmt.Println("Warning: failed to close response body:", err)
	}
}
func (p *JiraPlugin) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided")
	}

	switch args[0] {
	case "create-issue":
		if len(args) < 4 {
			return fmt.Errorf("usage: create-issue <projectKey> <summary> <description>")
		}
		return p.CreateIssue(args[1], args[2], args[3])
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

		issues, err := p.fetchAssignedIssues(args[1])
		if err != nil {
			return err
		}

		return p.generateAISummary(args[1], issues)
	default:
		return fmt.Errorf("unknown Jira command: %s", args[0])
	}
}

func (p *JiraPlugin) Info() (string, string) {
	return "jira", "Integration with Jira for tracking issues"
}

func (j *JiraPlugin) CreateIssue(projectKey, summary, description string) error {
	issue := map[string]interface{}{
		"fields": map[string]interface{}{
			"project":     map[string]string{"key": projectKey},
			"summary":     summary,
			"description": description,
			"issuetype":   map[string]string{"name": "Task"},
		},
	}

	body, err := toJSON(issue)
	if err != nil {
		return wrapError("failed to prepare issue payload", err)
	}

	resp, err := j.makeRequest("POST", "/rest/api/2/issue", bytes.NewReader(body), false)
	if err != nil {
		return wrapError("failed to create Jira issue", err)
	}

	defer HandleResponseBody(resp.Body)
	respBody, _ := io.ReadAll(resp.Body)
	logger.Logger.Debug().Str("JiraResponse", string(respBody)).Msg("Jira API response")
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected response from Jira: %s", resp.Status)
	}

	return nil
}

func (p *JiraPlugin) listIssues(projectKey string) error {
	endpoint := fmt.Sprintf("/rest/api/2/search?jql=project=\"%s\"", projectKey)

	resp, err := p.makeRequest("GET", endpoint, nil, false)
	if err != nil {
		return wrapError("failed to fetch Jira issues", err)
	}
	defer HandleResponseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to fetch issues, status: %s, response: %s", resp.Status, string(body))
	}

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

	fmt.Printf("\n📌 Issues for project **%s**:\n", projectKey)
	for _, issue := range result.Issues {
		fmt.Printf("   - [%s] %s\n", issue.Key, issue.Fields.Summary)
	}
	logger.Logger.Info().
		Str("Project", projectKey).
		Int("Issue Count", len(result.Issues)).
		Msg("Fetched Jira issues")

	return nil
}

func (p *JiraPlugin) assignedIssues(userEmail string) error {
	jql := fmt.Sprintf("assignee='%s' ORDER BY updated DESC", userEmail)
	encodedJQL := url.QueryEscape(jql)
	endpoint := fmt.Sprintf("/rest/api/2/search?jql=%s", encodedJQL)

	fmt.Println("Final Request URL:", p.baseURL+endpoint)

	resp, err := p.makeRequest("GET", endpoint, nil, false)
	if err != nil {
		return wrapError("failed to fetch assigned issues", err)
	}
	defer HandleResponseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to fetch assigned issues, status: %s, response: %s", resp.Status, string(body))
	}

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
		fmt.Printf("\n📌 No issues assigned to **%s**.\n", userEmail)
		return nil
	}

	fmt.Printf("\n📌 Issues assigned to **%s**:\n", userEmail)

	for _, issue := range result.Issues {
		fmt.Printf("   - [%s] (%s) %s\n", issue.Key, issue.Fields.IssueType.Name, issue.Fields.Summary)
		fmt.Printf("     🔹 Status: %s | 📅 Updated: %s\n", issue.Fields.Status.Name, issue.Fields.Updated)
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
	os.Getenv("OPENAI_ORG")

	if apiKey == "" {
		return fmt.Errorf("OPENAI_API_KEY is not set. Please export your API key.")
	}

	var formattedIssuesText strings.Builder
	for _, issue := range issues {
		formattedIssuesText.WriteString(fmt.Sprintf("- [%s] %s: %s (Status: %s, Updated: %s)\n",
			issue.IssueType, issue.Key, issue.Summary, issue.Status, issue.Updated))
	}

	prompt := fmt.Sprintf(`
      I am preparing a self-evaluation for my work. Please summarize my Jira contributions in a professional yet concise way.
      Focus on the impact of my work rather than just listing tasks. 
      Frame the summary as if I am personally describing my achievements for a performance review.

      Here are my recent Jira contributions:
      %s

      Respond in the first person, starting with "I...".
      Use natural language that sounds like something I would say in a self-assessment.
      Keep it concise and focused on impact.`, formattedIssuesText.String())

	payload := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "You are an assistant summarizing Jira issue contributions."},
			{"role": "user", "content": prompt},
		},
		"max_tokens": 200,
	}

	body, err := toJSON(payload)
	if err != nil {
		return wrapError("failed to prepare AI request payload", err)
	}

	resp, err := p.makeRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(body), true)
	if err != nil {
		return wrapError("failed to get AI summary", err)
	}
	defer HandleResponseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get AI summary, status: %s, response: %s", resp.Status, string(respBody))
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
	jql := fmt.Sprintf("assignee='%s' ORDER BY updated DESC", userEmail)
	encodedJQL := url.QueryEscape(jql)
	endpoint := fmt.Sprintf("/rest/api/2/search?jql=%s", encodedJQL)

	resp, err := p.makeRequest("GET", endpoint, nil, false)
	if err != nil {
		return nil, wrapError("failed to fetch assigned issues", err)
	}
	defer HandleResponseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch assigned issues, status: %s, response: %s", resp.Status, string(body))
	}

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

	formattedIssues := make([]struct {
		Key       string
		Summary   string
		IssueType string
		Status    string
		Updated   string
	}, len(result.Issues))

	for i, issue := range result.Issues {
		formattedIssues[i] = struct {
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
		}
	}

	return formattedIssues, nil
}
