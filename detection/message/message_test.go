package message

import (
	"testing"

	"github.com/chaoss/ai-detection-action/detection"
)

func TestDetect(t *testing.T) {
	d := &Detector{}
	tests := []struct {
		name      string
		message   string
		wantTools []string
	}{
		{
			name:      "aider prefix",
			message:   "aider: fix the login bug",
			wantTools: []string{"Aider"},
		},
		{
			name:      "aider prefix uppercase",
			message:   "Aider: refactor auth module",
			wantTools: []string{"Aider"},
		},
		{
			name:      "Claude Code footer",
			message:   "Add user validation\n\nGenerated with Claude Code",
			wantTools: []string{"Claude Code"},
		},
		{
			name:      "Claude Code footer with link",
			message:   "Add validation\n\nGenerated with Claude Code\nhttps://claude.ai",
			wantTools: []string{"Claude Code"},
		},
		{
			name:      "no patterns",
			message:   "normal commit message with no AI signatures",
			wantTools: nil,
		},
		{
			name:      "aider in middle of message not prefix",
			message:   "fix the aider: integration test",
			wantTools: nil,
		},
		{
			name:      "aider as substring of a word",
			message:   "raider: fix the tests",
			wantTools: nil,
		},
		{
			name:      "empty message",
			message:   "",
			wantTools: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := d.Detect(detection.Input{CommitMessage: tt.message})
			gotTools := make([]string, len(findings))
			for i, f := range findings {
				gotTools[i] = f.Tool
				if f.Confidence != detection.ConfidenceMedium {
					t.Errorf("confidence = %d, want %d", f.Confidence, detection.ConfidenceMedium)
				}
				if f.Detector != "message" {
					t.Errorf("detector = %q, want %q", f.Detector, "message")
				}
			}

			if len(gotTools) == 0 {
				gotTools = nil
			}

			if len(gotTools) != len(tt.wantTools) {
				t.Errorf("tools = %v, want %v", gotTools, tt.wantTools)
				return
			}
			for i := range gotTools {
				if gotTools[i] != tt.wantTools[i] {
					t.Errorf("tools = %v, want %v", gotTools, tt.wantTools)
					return
				}
			}
		})
	}
}
