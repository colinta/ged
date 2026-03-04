package rule

import (
	"bytes"
	"os/exec"
	"strings"
)

// XargsRule executes a shell command for each line, appending the line as an argument.
// The command's stdout replaces the line. If the command fails, the line passes through unchanged.
type XargsRule struct {
	command string
}

// NewXargsRule creates a new XargsRule.
func NewXargsRule(command string) *XargsRule {
	return &XargsRule{command: command}
}

// Apply runs the command with the line appended as an argument.
func (r *XargsRule) Apply(line string, ctx *LineContext) ([]string, error) {
	// Build the full shell command: "command line_text"
	fullCmd := r.command + " " + shellQuote(line)
	cmd := exec.Command("sh", "-c", fullCmd)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		// Command failed — pass line through unchanged
		return []string{line}, nil
	}

	// Split output into lines, trimming trailing newline
	output := strings.TrimRight(stdout.String(), "\n")
	if output == "" {
		return []string{}, nil
	}
	return strings.Split(output, "\n"), nil
}

// shellQuote wraps a string in single quotes for safe shell use.
// Single quotes inside the string are escaped as '\'' (end quote, literal quote, start quote).
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
