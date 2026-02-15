package main

import (
	"regexp"
	"strings"
)

var knownAgentCommitters = map[string]string{
	"209825114+claude[bot]@users.noreply.github.com":					"Claude",
	"215619710+anthropic-claude[bot]@users.noreply.github.com":			"Claude (Anthropic)",
	"208546643+claude-code-action[bot]@users.noreply.github.com": 		"Claude Code Action",
	"198982749+copilot@users.noreply.github.com": 						"GitHub Copilot (agent)",
	"167198135+copilot[bot]@users.noreply.github.com": 					"GitHub Copilot (chat)",
	"206951365+cursor[bot]@users.noreply.github.com": 					"Cursor",
	"215057067+openai-codex[bot]@users.noreply.github.com": 			"OpenAI Codex",
	"199175422+chatgpt-codex-connector[bot]@users.noreply.github.com": 	"Codex via ChatGPT",
	"176961590+gemini-code-assist[bot]@users.noreply.github.com": 		"Gemini Code Assist",
	"208079219+amazon-q-developer[bot]@users.noreply.github.com": 		"Amazon Q Developer",
	"158243242+devin-ai-integration[bot]@users.noreply.github.com": 	"Devin",
	"205137888+cline[bot]@users.noreply.github.com": 					"Cline",
	"230936708+continue[bot]@users.noreply.github.com": 				"Continue.dev",
	"201248094+sourcegraph-cody[bot]@users.noreply.github.com": 		"Sourcegraph Cody",
	"220155983+jetbrains-ai[bot]@users.noreply.github.com": 			"JetBrains AI",
	"136622811+coderabbitai[bot]@users.noreply.github.com": 			"CodeRabbit",
}

// Check if a given committer email is a known AI agent, returns the agent name and a boolean for success
func IdentifyAgentCommitterEmail(email string) (string, bool) {
	cleanEmail := strings.ToLower(strings.TrimSpace(email))

	name, exists := knownAgentCommitters[cleanEmail]
	return name, exists
}

var knownCoAuthorEmails = map[string]string{
	"noreply@anthropic.com": "Claude Code",
	"cursoragent@cursor.com": "Cursor",
	"noreply@aider.chat":    "Aider",
}

var coAuthorPattern = regexp.MustCompile(`(?im)^co-authored-by:\s*[^<]*<([^>]+)>`)

// IdentifyCoAuthorTrailers parses a commit message for Co-Authored-By trailers
// and returns the names of any matched AI tools.
func IdentifyCoAuthorTrailers(message string) []string {
	matches := coAuthorPattern.FindAllStringSubmatch(message, -1)
	var tools []string
	seen := map[string]bool{}
	for _, match := range matches {
		email := strings.ToLower(strings.TrimSpace(match[1]))
		if name, ok := knownCoAuthorEmails[email]; ok && !seen[name] {
			tools = append(tools, name)
			seen[name] = true
		}
	}
	return tools
}

var commitMessagePatterns = []struct {
	check func(string) bool
	name  string
}{
	{
		check: func(msg string) bool {
			return strings.HasPrefix(strings.ToLower(msg), "aider:")
		},
		name: "Aider",
	},
	{
		check: func(msg string) bool {
			return strings.Contains(msg, "Generated with Claude Code")
		},
		name: "Claude Code",
	},
}

// IdentifyCommitMessagePatterns checks a commit message for known AI tool
// signatures and returns matched tool names.
func IdentifyCommitMessagePatterns(message string) []string {
	var tools []string
	for _, p := range commitMessagePatterns {
		if p.check(message) {
			tools = append(tools, p.name)
		}
	}
	return tools
}