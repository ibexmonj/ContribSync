# 🚀 ContribSync (`csync`) - Track & Summarize Your Contributions


⚠️ **This is an early-stage open-source project.**


ContribSync (`csync`) is a CLI tool that helps individual contributors track their work across multiple platforms, generate AI-powered summaries, and streamline performance reviews.

## 📌 Features

### ✅ Jira Integration
- **Create issues** in Jira directly from the CLI.
- **List issues** in a given Jira project.
- **Find issues assigned to a user.**
- **Generate AI-powered summaries** of Jira contributions for self-evaluations (using OpenAI GPT-3.5).

### ✅ Pluggable Architecture
- Built to **support multiple integrations** (e.g., GitHub, GitLab, Slack, Teams in the future).

### ✅ Automated Summarization
- Converts Jira activity into a **structured, self-evaluation-friendly summary**.
- Helps users **document their contributions efficiently** for performance reviews.


## 🛠️ Setup & Running Locally

### Prerequisites
- Go **1.22.4+** installed
- Jira API access (API token required)
- (Optional) OpenAI API key for AI-generated summaries

### Install & Build
Clone the repo and navigate into the project directory:
```sh
git clone https://github.com/your-org/ContribSync.git
cd ContribSync
go build -o csync main.go
```

### Run the CLI
```sh
./csync --help
```

### 🌍 Environment Variables

Before running commands, ensure the required environment variables are set.

✅ Jira API Configuration
```sh
export JIRA_BASE_URL=https://your-org.atlassian.net
export JIRA_API_TOKEN=your-api-token
export JIRA_EMAIL="your-email@example.com"
```

✅ (Optional) OpenAI API Configuration
```sh
export OPENAI_ORG=your-org-id
```
🏗️ Example Commands

📌 Create a Jira Issue
```sh  
./csync plugin exec jira create-issue "Optimize Sync Performance" "Improve background sync to reduce CPU usage."
```
Sample Output:
``` 
Issue created successfully:
Key: CSYNC-101
Summary: Optimize Sync Performance
```

📌 List Issues in a Project
```sh
./csync plugin exec jira list-issues PROJECT_KEY //project key under JIRA project settings
```
Sample Output:
```
Issues in PROJECT_KEY:
- CSYNC-101: Optimize Sync Performance (To Do)
- CSYNC-102: Implement API Rate Limiting (In Progress)
- CSYNC-103: Fix Sync Conflict Errors (Done)
```

📌 Find Issues Assigned to a User
```sh
./csync plugin exec jira assigned-issues your-email@example.com
```
Sample Output:
```
📌 Issues assigned to your-email@example.com:
- CSYNC-101: Optimize Sync Performance (To Do)
- CSYNC-102: Implement API Rate Limiting (In Progress)
- CSYNC-103: Fix Sync Conflict Errors (Done)
```

📌 Generate AI-Powered Summary for Self-Evaluation
```sh
./csync plugin exec jira summary your-email@example.com
```
Sample Output:
```
📌 AI-Generated Self-Evaluation Summary:
Over the past cycle, I worked on key improvements to our CloudSync platform:
- Resolved **file synchronization conflicts**, preventing duplicate file issues (CSYNC-103).
- Implemented **API rate limiting**, improving system stability under heavy load (CSYNC-102).
- Optimized **background sync performance**, reducing CPU usage and increasing efficiency (CSYNC-101).
These contributions enhanced platform reliability and performance, benefiting both end users and internal teams.
```
_This requires the OpenAI API key to be set._

### ✅ GitHub Integration
- **Fetch pull requests** from a repository.
- **List commits** associated with each PR.
- **Track merged and closed PRs** for contribution logging.

_Requires a valid GITHUB_TOKEN to be set._
```sh
export GITHUB_TOKEN=your-personal-access-token
```


📌 Fetch Pull Requests & Commits from GitHub
```sh
./csync plugin exec github summary owner/repo

```

Sample Output:
```

📌 Pull Request Summary for owner/repo

🔹 PR #42: Fix database connection issue (closed)
🏷️ Status: closed
🔄 Merged: true
📆 Created: 2024-06-15 00:49:33 UTC
📝 Commits:
- [abc123] Fix DB connection timeout
- [def456] Improve error logging
```

## 🚀 We’re Adding Features Regularly!

This project is evolving, and we’re actively adding new integrations and improvements.
Watch this space for updates! 🚀

## 📜 License

This project is licensed under the [MIT License.](LICENSE)

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute.

## 🛠️ Future Roadmap
•	Refactor cli command structure . E.g. "csync plugin exec jira..." to "csync jira..."
•	Slack & Teams Notifications  
•	Export Summaries as JSON or Markdown  

## 🌟 Support & Feedback

Have ideas or feedback? Open an issue or reach out!