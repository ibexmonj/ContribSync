# 🚀 ContribSync (`csync`) - Track & Summarize Your Contributions

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

📌 List Issues in a Project
```sh
./csync plugin exec jira list-issues PROJECT_KEY //project key under JIRA project settings
```
📌 Find Issues Assigned to a User
```sh
./csync plugin exec jira assigned-issues your-email@example.com
```
📌 Generate AI-Powered Summary for Self-Evaluation
```sh
./csync plugin exec jira summary your-email@example.com
```
_This requires the OpenAI API key to be set._


## 🚀 We’re Adding Features Regularly!

This project is evolving, and we’re actively adding new integrations and improvements.
Watch this space for updates!
Want to suggest a feature? Open an issue or a pull request! 🚀

## 📜 License

This project is licensed under the [MIT License.](LICENSE)

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute.

## 🛠️ Future Roadmap
•	GitHub Integration (track PRs, commits, issues)
•	Slack & Teams Notifications
•	Export Summaries as JSON or Markdown

## 🌟 Support & Feedback

Have ideas or feedback? Open an issue or reach out!