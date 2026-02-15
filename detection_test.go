package main

import (
	"slices"
	"testing"
)

func TestIdentifyAgentCommitterEmailAllKnown(t *testing.T) {
	for email, expectedName := range knownAgentCommitters {
		name, exists := IdentifyAgentCommitterEmail(email)
		if !exists || name != expectedName {
			t.Errorf("IdentifyAgentCommitterEmail(%q) = %q, %v, want %q, true", email, name, exists, expectedName)
		}
	}
}

func TestIdentifyAgentCommitterEmailMixedCase(t *testing.T) {
	cases := []struct {
		input    string
		wantName string
	}{
		{"198982749+Copilot@users.noreply.github.com", "GitHub Copilot (agent)"},
		{"209825114+CLAUDE[BOT]@USERS.NOREPLY.GITHUB.COM", "Claude"},
		{"136622811+CodeRabbitAI[bot]@users.noreply.github.com", "CodeRabbit"},
	}

	for _, tc := range cases {
		name, exists := IdentifyAgentCommitterEmail(tc.input)
		if !exists || name != tc.wantName {
			t.Errorf("IdentifyAgentCommitterEmail(%q) = %q, %v, want %q, true", tc.input, name, exists, tc.wantName)
		}
	}
}

func TestIdentifyAgentCommitterEmailWhitespace(t *testing.T) {
	cases := []string{
		"  209825114+claude[bot]@users.noreply.github.com",
		"209825114+claude[bot]@users.noreply.github.com  ",
		"  209825114+claude[bot]@users.noreply.github.com  ",
	}

	for _, email := range cases {
		name, exists := IdentifyAgentCommitterEmail(email)
		if !exists || name != "Claude" {
			t.Errorf("IdentifyAgentCommitterEmail(%q) = %q, %v, want %q, true", email, name, exists, "Claude")
		}
	}
}

func TestIdentifyAgentCommitterEmailNotFound(t *testing.T) {
	cases := []string{
		"user@example.com",
		"",
		"  ",
		"claude@anthropic.com",
		"not-a-bot@users.noreply.github.com",
	}

	for _, email := range cases {
		name, exists := IdentifyAgentCommitterEmail(email)
		if exists || name != "" {
			t.Errorf("IdentifyAgentCommitterEmail(%q) = %q, %v, want \"\", false", email, name, exists)
		}
	}
}

func TestIdentifyCoAuthorTrailers(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    []string
	}{
		{
			name:    "Claude trailer with Opus model",
			message: "fix: update handler\n\nCo-Authored-By: Claude Opus 4 <noreply@anthropic.com>",
			want:    []string{"Claude Code"},
		},
		{
			name:    "Claude trailer with Sonnet model",
			message: "fix: update handler\n\nCo-Authored-By: Claude Sonnet 4 <noreply@anthropic.com>",
			want:    []string{"Claude Code"},
		},
		{
			name:    "Cursor trailer",
			message: "refactor: extract method\n\nCo-Authored-By: Cursor <cursoragent@cursor.com>",
			want:    []string{"Cursor"},
		},
		{
			name:    "Aider trailer with model name",
			message: "feat: add endpoint\n\nCo-Authored-By: aider (gpt-4o) <noreply@aider.chat>",
			want:    []string{"Aider"},
		},
		{
			name:    "Aider trailer with different model",
			message: "feat: add endpoint\n\nCo-Authored-By: aider (claude-3.5-sonnet) <noreply@aider.chat>",
			want:    []string{"Aider"},
		},
		{
			name:    "multiple trailers with Claude and human",
			message: "fix: bug\n\nCo-Authored-By: Claude Opus 4 <noreply@anthropic.com>\nCo-Authored-By: Alice <alice@example.com>",
			want:    []string{"Claude Code"},
		},
		{
			name:    "multiple AI trailers",
			message: "fix: bug\n\nCo-Authored-By: Claude Opus 4 <noreply@anthropic.com>\nCo-Authored-By: aider (gpt-4o) <noreply@aider.chat>",
			want:    []string{"Claude Code", "Aider"},
		},
		{
			name:    "case variation",
			message: "fix: thing\n\nco-authored-by: Claude <noreply@anthropic.com>",
			want:    []string{"Claude Code"},
		},
		{
			name:    "CO-AUTHORED-BY uppercase",
			message: "fix: thing\n\nCO-AUTHORED-BY: Claude <noreply@anthropic.com>",
			want:    []string{"Claude Code"},
		},
		{
			name:    "no trailers",
			message: "just a normal commit message",
			want:    nil,
		},
		{
			name:    "human co-author only",
			message: "pair programming\n\nCo-Authored-By: Bob <bob@company.com>",
			want:    nil,
		},
		{
			name:    "empty message",
			message: "",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IdentifyCoAuthorTrailers(tt.message)
			if !slices.Equal(got, tt.want) {
				t.Errorf("IdentifyCoAuthorTrailers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentifyCommitMessagePatterns(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    []string
	}{
		{
			name:    "aider prefix",
			message: "aider: fix the login bug",
			want:    []string{"Aider"},
		},
		{
			name:    "aider prefix uppercase",
			message: "Aider: refactor auth module",
			want:    []string{"Aider"},
		},
		{
			name:    "Claude Code footer",
			message: "Add user validation\n\nGenerated with Claude Code",
			want:    []string{"Claude Code"},
		},
		{
			name:    "Claude Code footer with link",
			message: "Add validation\n\nGenerated with Claude Code\nhttps://claude.ai",
			want:    []string{"Claude Code"},
		},
		{
			name:    "no patterns",
			message: "normal commit message with no AI signatures",
			want:    nil,
		},
		{
			name:    "aider in middle of message not prefix",
			message: "fix the aider: integration test",
			want:    nil,
		},
		{
			name:    "aider as substring of a word",
			message: "raider: fix the tests",
			want:    nil,
		},
		{
			name:    "empty message",
			message: "",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IdentifyCommitMessagePatterns(tt.message)
			if !slices.Equal(got, tt.want) {
				t.Errorf("IdentifyCommitMessagePatterns() = %v, want %v", got, tt.want)
			}
		})
	}
}

