package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/chaoss/ai-detection-action/detection"
	"github.com/chaoss/ai-detection-action/detection/coauthor"
	"github.com/chaoss/ai-detection-action/detection/committer"
	"github.com/chaoss/ai-detection-action/detection/message"
	"github.com/chaoss/ai-detection-action/detection/toolmention"
	"github.com/chaoss/ai-detection-action/output"
	"github.com/chaoss/ai-detection-action/scan"
)

var Version = "dev"

// Exit codes
const (
	ExitNoAI    = 0
	ExitAI      = 1
	ExitError   = 2
)

func allDetectors() []detection.Detector {
	return []detection.Detector{
		&committer.Detector{},
		&coauthor.Detector{},
		&message.Detector{},
		&toolmention.Detector{},
	}
}

// Run is the main entry point for the CLI. Returns an exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(stderr, "usage: ai-detection <command> [options]")
		fmt.Fprintln(stderr, "commands: commits, text, version")
		return ExitError
	}

	switch args[0] {
	case "commits":
		return runCommits(args[1:], stdout, stderr)
	case "text":
		return runText(args[1:], stdout, stderr)
	case "version":
		fmt.Fprintf(stdout, "ai-detection %s\n", Version)
		return ExitNoAI
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		fmt.Fprintln(stderr, "commands: commits, text, version")
		return ExitError
	}
}

func runCommits(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("commits", flag.ContinueOnError)
	fs.SetOutput(stderr)

	rangeFlag := fs.String("range", "", "commit range in BASE..HEAD format")
	formatFlag := fs.String("format", "text", "output format: json or text")
	minConfFlag := fs.String("min-confidence", "low", "minimum confidence level: low, medium, high (or 1, 2, 3)")

	if err := fs.Parse(args); err != nil {
		return ExitError
	}

	repoPath := "."
	if fs.NArg() > 0 {
		repoPath = fs.Arg(0)
	}

	minConf, err := output.ConfidenceFromString(*minConfFlag)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return ExitError
	}

	detectors := allDetectors()
	report, err := scan.ScanCommitRange(repoPath, *rangeFlag, detectors)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return ExitError
	}

	report = filterReport(report, minConf)

	switch *formatFlag {
	case "json":
		if err := output.FormatJSON(stdout, report); err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return ExitError
		}
	case "text":
		if err := output.FormatText(stdout, report); err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return ExitError
		}
	default:
		fmt.Fprintf(stderr, "unknown format: %s\n", *formatFlag)
		return ExitError
	}

	if report.Summary.AICommits > 0 {
		return ExitAI
	}
	return ExitNoAI
}

func runText(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("text", flag.ContinueOnError)
	fs.SetOutput(stderr)

	formatFlag := fs.String("format", "text", "output format: json or text")
	inputFlag := fs.String("input", "-", "input file path, or - for stdin")

	if err := fs.Parse(args); err != nil {
		return ExitError
	}

	var textBytes []byte
	var err error

	if *inputFlag == "-" {
		textBytes, err = io.ReadAll(os.Stdin)
	} else {
		textBytes, err = os.ReadFile(*inputFlag)
	}
	if err != nil {
		fmt.Fprintf(stderr, "error reading input: %v\n", err)
		return ExitError
	}

	detectors := allDetectors()
	findings := scan.ScanText(string(textBytes), detectors)

	switch *formatFlag {
	case "json":
		if err := output.FormatJSONFindings(stdout, findings); err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return ExitError
		}
	case "text":
		if err := output.FormatTextFindings(stdout, findings); err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return ExitError
		}
	default:
		fmt.Fprintf(stderr, "unknown format: %s\n", *formatFlag)
		return ExitError
	}

	if len(findings) > 0 {
		return ExitAI
	}
	return ExitNoAI
}

func filterReport(report scan.Report, minConf detection.Confidence) scan.Report {
	if minConf <= detection.ConfidenceLow {
		return report
	}

	filtered := scan.Report{
		Commits: make([]scan.CommitResult, 0, len(report.Commits)),
		Summary: scan.Summary{
			TotalCommits: report.Summary.TotalCommits,
			ToolCounts:   map[string]int{},
			ByConfidence: map[string]int{},
		},
	}

	for _, cr := range report.Commits {
		var kept []detection.Finding
		for _, f := range cr.Findings {
			if f.Confidence >= minConf {
				kept = append(kept, f)
			}
		}
		result := scan.CommitResult{Hash: cr.Hash, Findings: kept}
		filtered.Commits = append(filtered.Commits, result)

		if len(kept) > 0 {
			filtered.Summary.AICommits++
		}
		for _, f := range kept {
			filtered.Summary.ToolCounts[f.Tool]++
			filtered.Summary.ByConfidence[f.Confidence.String()]++
		}
	}

	return filtered
}
