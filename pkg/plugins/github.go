package plugins

import (
	"context"
	"errors"
	"fmt"
	"github.com/ibexmonj/ContribSync/pkg/logger"
	"os"
	"strings"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

type GitHubPlugin struct{}

func (g *GitHubPlugin) Init() error {
	logger.Logger.Info().Msg("✅ GitHub plugin initialized")
	return nil
}

func (g *GitHubPlugin) Info() (string, string) {
	return "github", "GitHub Plugin: Fetch PRs and commits"
}

func (g *GitHubPlugin) Execute(args []string) error {
	if len(args) < 2 {
		return errors.New("Usage: csync plugin exec github summary owner/repo [email]")
	}

	if args[0] != "summary" {
		return fmt.Errorf("Unknown command for github: %s", args[0])
	}

	repoArg := args[1]
	owner, repo, err := parseOwnerRepo(repoArg)
	if err != nil {
		return err
	}

	var emailFilter string

	if len(args) >= 3 {
		emailFilter = args[2]
	}

	return GitHubSummary(owner, repo, emailFilter)
}

// GitHubSummary fetches PRs & commits for a repo
func GitHubSummary(owner, repo, emailFilter string) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.New("❌ GITHUB_TOKEN is not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	prs, err := fetchPRs(client, ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("❌ Failed to fetch PRs: %w", err)
	}

	if emailFilter != "" {
		logger.Logger.Info().
			Str("owner", owner).
			Str("repo", repo).
			Msg("📌 Pull Request Summary")
	} else {
		logger.Logger.Info().
			Str("owner", owner).
			Str("repo", repo).
			Msg("📌 Pull Request Summary")
	}
	prCount := 0

	for _, pr := range prs {

		commits, err := fetchCommits(client, ctx, owner, repo, *pr.Number)
		if err != nil {
			fmt.Printf("   ⚠️ Failed to fetch commits for PR #%d: %v\n", *pr.Number, err)
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

		fmt.Printf("\n🔹 **PR #%d**: %s (%s)\n", *pr.Number, *pr.Title, *pr.State)
		fmt.Printf("   🏷️ Status: %s | 🔄 Merged: %v | 📆 Created: %v\n", *pr.State, pr.MergedAt != nil, pr.CreatedAt)

		if len(matchingCommits) > 0 {
			fmt.Println("   📝 Commits:")
			for _, commit := range matchingCommits {
				fmt.Printf("      - [%s] %s\n", (*commit.SHA)[:7], *commit.Commit.Message)
			}
		} else {
			fmt.Println("   🚫 No matching commits found.")
		}
	}

	if prCount == 0 {
		fmt.Println("\n❌ No pull requests found.")
	}
	return nil
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
	if err != nil {
		return nil, fmt.Errorf("error fetching PRs from GitHub: %w", err)
	}
	return prs, nil
}

// Fetch commits for a specific PR
func fetchCommits(client *github.Client, ctx context.Context, owner, repo string, prNumber int) ([]*github.RepositoryCommit, error) {
	commits, _, err := client.PullRequests.ListCommits(ctx, owner, repo, prNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching commits for PR #%d: %w", prNumber, err)
	}
	return commits, nil
}

// Parse "owner/repo" format
func parseOwnerRepo(full string) (string, string, error) {
	parts := strings.Split(full, "/")
	if len(parts) != 2 {
		return "", "", errors.New("invalid repo format, expected owner/repo")
	}
	return parts[0], parts[1], nil
}
