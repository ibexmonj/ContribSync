# 🚀 ContribSync (`csync`) - Track & Summarize Your Contributions


⚠️ **This is an early-stage open-source project.**

> ⚡️ Stop dreading performance reviews. Start tracking your work effortlessly.

**ContribSync** is a CLI tool that automatically fetches your contributions from Jira, GitHub, and more—then generates a clean, AI-powered self-evaluation draft.  
It's built for engineers who want to spend 2 minutes, not 2 days, preparing for performance reviews.

---

## ✨ Why Use ContribSync?

Every quarter, we scramble to remember what we actually shipped. Hunting through Jira tickets, pull requests, and Slack threads wastes time and leads to incomplete self-reviews.

**ContribSync automates that painful process.**  
It connects to the tools you're already using and summarizes what you’ve done—so you can focus on impact, not busywork.

---


## 📌 Features

### ✅ Jira Integration
- Fetches assigned tickets from Jira (JQL-based)
- Summarizes tickets using GPT-3.5 / GPT-4
- Outputs clean Markdown for performance reviews

### ✅ GitHub Plugin
- Lists PRs created or merged in a timeframe
- Extracts commit messages and PR descriptions
- Generates activity summaries via AI

### ✅ Slack Plugin (WIP)
- Sends daily or weekly contribution reminder notifications
- Future: Slack-based input capture for non-CLI users

### ✅ AI-Powered Summaries
- Uses OpenAI to turn raw issue/PR data into concise summaries
- Fully configurable prompt structure
- Offline/manual mode available for auditability


## 🛠️ Setup & Running Locally

### Prerequisites
- Go **1.22.4+** installed
- Jira API access (API token required)
- (Optional) OpenAI API key for AI-generated summaries

### Install & Build
Clone the repo and navigate into the project directory:
```sh
git clone https://github.com/your-org/ContribSync.git
cd ContribSync/cmd/csync
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

# Optional: GitHub
GITHUB_TOKEN=ghp_...

# Optional: Slack
SLACK_WEBHOOK_URL=https://hooks.slack.com/...

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
- **Track  open, merged and closed PRs** for contribution logging.

_Requires a valid GITHUB_TOKEN to be set._
```sh
export GITHUB_TOKEN=your-personal-access-token
```


📌 Fetch Pull Requests & Commits from GitHub; filter by email
```sh
./csync plugin exec github summary owner/repo email@example.com

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

## 🧩 Plugin System

ContribSync is built to be extensible. Add your own integrations via plugins under pkg/plugins.

Available plugins:
•	Jira  
•	GitHub  
•	Slack (WIP)  

## 🧠 Powered by OpenAI

ContribSync uses the OpenAI API to generate human-readable summaries. You control the prompt format and engine (gpt-3.5-turbo or gpt-4).

## 📜 License

This project is licensed under the [MIT License.](LICENSE)

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute.

We’re actively building this and looking for contributors!

Ways to contribute:  
•	Add a new platform plugin (Linear, GitLab, Trello)  
•	Improve AI summarization prompts  
•	Build a simple TUI/GUI wrapper  
•	Add caching / config improvements  
•	Write tests and docs  

## 🛠️ Future Roadmap
•	Refactor cli command structure . E.g. "csync plugin exec jira..." to "csync jira..."
•	Slack & Teams Notifications  
•	Export Summaries as JSON or Markdown  

## 🌟 Support & Feedback

Have ideas or feedback? Open an issue or reach out!