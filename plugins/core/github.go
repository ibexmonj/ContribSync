package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// GitHubPlugin implements the Plugin interface
type GitHubPlugin struct{}

// Init initializes the GitHub plugin
func (g *GitHubPlugin) Init() error {
	fmt.Println("✅ GitHub plugin initialized")
	return nil
}

// Info returns plugin metadata
func (g *GitHubPlugin) Info() (string, string) {
	return "github", "GitHub Plugin: Fetch PRs and commits"
}

// Execute runs the GitHub plugin with given args
func (g *GitHubPlugin) Execute(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Usage: csync plugin exec github summary owner/repo")
	}

	if args[0] != "summary" {
		return fmt.Errorf("Unknown command for github: %s", args[0])
	}

	repoArg := args[1]
	owner, repo := parseOwnerRepo(repoArg)

	GitHubSummary(owner, repo)
	return nil
}

// GitHubSummary fetches PRs & commits for a repo
func GitHubSummary(owner, repo string) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	prs, err := fetchPRs(client, ctx, owner, repo)
	if err != nil {
		log.Fatalf("Error fetching PRs: %v", err)
	}

	fmt.Printf("📌 Pull Request Summary for %s/%s\n", owner, repo)
	for _, pr := range prs {
		fmt.Printf("\n🔹 PR #%d: %s (%s)\n", *pr.Number, *pr.Title, *pr.State)
		fmt.Printf("   🏷️ Status: %s\n", *pr.State)
		fmt.Printf("   🔄 Merged: %v\n", pr.MergedAt != nil)
		fmt.Printf("   📆 Created: %v\n", pr.CreatedAt)

		commits, err := fetchCommits(client, ctx, owner, repo, *pr.Number)
		if err != nil {
			fmt.Printf("   ⚠️ Error fetching commits: %v\n", err)
			continue
		}

		fmt.Printf("   📝 Commits:\n")
		for _, commit := range commits {
			fmt.Printf("      - [%s] %s\n", (*commit.SHA)[:7], *commit.Commit.Message)
		}
	}
}

// Fetch PRs for the repository
func fetchPRs(client *github.Client, ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
	opts := &github.PullRequestListOptions{State: "all", ListOptions: github.ListOptions{PerPage: 10}}
	prs, _, err := client.PullRequests.List(ctx, owner, repo, opts)
	return prs, err
}

// Fetch commits for a specific PR
func fetchCommits(client *github.Client, ctx context.Context, owner, repo string, prNumber int) ([]*github.RepositoryCommit, error) {
	commits, _, err := client.PullRequests.ListCommits(ctx, owner, repo, prNumber, nil)
	return commits, err
}

// Parse "owner/repo" format
func parseOwnerRepo(full string) (string, string) {
	parts := strings.Split(full, "/")
	if len(parts) != 2 {
		log.Fatal("Invalid repo format, expected owner/repo")
	}
	return parts[0], parts[1]
}
