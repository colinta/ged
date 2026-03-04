package rule

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ExecRule pipes the entire document through a shell command.
// The command receives the document on stdin and its stdout becomes the new document.
type ExecRule struct {
	command string
}

// NewExecRule creates a new ExecRule.
func NewExecRule(command string) *ExecRule {
	return &ExecRule{command: command}
}

// ApplyDocument pipes all lines through the command.
func (r *ExecRule) ApplyDocument(lines []string) ([]string, error) {
	cmd := exec.Command("sh", "-c", r.command)

	// Feed the document as stdin
	cmd.Stdin = strings.NewReader(strings.Join(lines, "\n") + "\n")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("exec %q failed: %w\n%s", r.command, err, stderr.String())
	}

	// Split output into lines, trimming trailing newline
	output := strings.TrimRight(stdout.String(), "\n")
	if output == "" {
		return []string{}, nil
	}
	return strings.Split(output, "\n"), nil
}
