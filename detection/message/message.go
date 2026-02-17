package message

import (
	"fmt"
	"strings"

	"github.com/chaoss/ai-detection-action/detection"
)

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

type Detector struct{}

func (d *Detector) Name() string { return "message" }

func (d *Detector) Detect(input detection.Input) []detection.Finding {
	if input.CommitMessage == "" {
		return nil
	}

	var findings []detection.Finding
	for _, p := range commitMessagePatterns {
		if p.check(input.CommitMessage) {
			findings = append(findings, detection.Finding{
				Detector:   d.Name(),
				Tool:       p.name,
				Confidence: detection.ConfidenceMedium,
				Detail:     fmt.Sprintf("commit message matches %s pattern", p.name),
			})
		}
	}

	return findings
}
