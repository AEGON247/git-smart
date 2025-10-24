// internal/git/git.go
package git

import (
	"errors"
	"os/exec"
	"strings"
)

// Run runs a git command and returns its trimmed output or an error.
// It combines stdout and stderr.
func Run(args ...string) (string, error) {
	// We're running the 'git' command. The args are passed in.
	cmd := exec.Command("git", args...)

	// CombinedOutput runs the command and returns its combined
	// standard output and standard error.
	output, err := cmd.CombinedOutput()

	// Trim whitespace from the output
	trimmedOutput := strings.TrimSpace(string(output))

	if err != nil {
		// If git command fails, return the error output
		return trimmedOutput, errors.New(trimmedOutput)
	}

	return trimmedOutput, nil
}

// IsRepo checks if the current directory is a git repository.
func IsRepo() bool {
	_, err := Run("rev-parse", "--is-inside-work-tree")
	return err == nil
}

// GetCurrentBranch returns the name of the current branch.
func GetCurrentBranch() (string, error) {
	return Run("rev-parse", "--abbrev-ref", "HEAD")
}
