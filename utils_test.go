package main

import (
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
