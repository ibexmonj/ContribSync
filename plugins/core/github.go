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
		return fmt.Errorf("Usage: csync plugin exec github summary owner/repo [email]")
	}

	if args[0] != "summary" {
		return fmt.Errorf("Unknown command for github: %s", args[0])
	}

	repoArg := args[1]
	owner, repo := parseOwnerRepo(repoArg)

	var emailFilter string

	if len(args) >= 3 {
		emailFilter = args[2]
	}

	GitHubSummary(owner, repo, emailFilter)
	return nil
}

// GitHubSummary fetches PRs & commits for a repo
func GitHubSummary(owner, repo, emailFilter string) {
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

	if emailFilter != "" {
		fmt.Printf("📌 Pull Request Summary for %s/%s (Filtered by commits from: %s)\n", owner, repo, emailFilter)
	} else {
		fmt.Printf("📌 Pull Request Summary for %s/%s\n", owner, repo)
	}
	prCount := 0

	for _, pr := range prs {

		commits, err := fetchCommits(client, ctx, owner, repo, *pr.Number)
		if err != nil {
			fmt.Printf("   ⚠️ Error fetching commits: %v\n", err)
			continue
		}

		var matchingCommits []*github.RepositoryCommit
		if emailFilter != "" {
			matchingCommits = filterCommitsByEmail(commits, emailFilter)
			if len(matchingCommits) == 0 {
				continue
			}
		} else {
			matchingCommits = commits
		}

		prCount++

		fmt.Printf("\n🔹 PR #%d: %s (%s)\n", *pr.Number, *pr.Title, *pr.State)
		fmt.Printf("   🏷️ Status: %s\n", *pr.State)
		fmt.Printf("   🔄 Merged: %v\n", pr.MergedAt != nil)
		fmt.Printf("   📆 Created: %v\n", pr.CreatedAt)

		fmt.Printf("   📝 Commits:\n")
		for _, commit := range matchingCommits {
			fmt.Printf("      - [%s] %s\n", (*commit.SHA)[:7], *commit.Commit.Message)
		}
	}
	if prCount == 0 && emailFilter != "" {
		fmt.Printf("\n❌ No pull requests found with commits by %s\n", emailFilter)
	}
}

func filterCommitsByEmail(commits []*github.RepositoryCommit, email string) []*github.RepositoryCommit {
	var filtered []*github.RepositoryCommit
	for _, commit := range commits {
		if commit.Commit.Author != nil && commit.Commit.Author.Email != nil && *commit.Commit.Author.Email == email {
			filtered = append(filtered, commit)
		}
	}
	return filtered
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
