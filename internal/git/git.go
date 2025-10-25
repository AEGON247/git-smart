// internal/git/git.go
package git

import (
	"errors"
	"os/exec"
	"regexp" // <-- Import regexp
	"strings"
)

// Run runs a git command and returns its trimmed output or an error.
// It combines stdout and stderr.
func Run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	trimmedOutput := strings.TrimSpace(string(output))

	if err != nil {
		// If git command fails, return the error output as the error message
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

// GetDefaultBranch finds the default branch name for the 'origin' remote.
func GetDefaultBranch() (string, error) {
	// Run the command `git remote show origin`
	output, err := Run("remote", "show", "origin")
	if err != nil {
		// Try 'master' as a fallback for purely local repos? Or just error out?
		// For now, let's error out if 'origin' isn't set up.
		return "", errors.New("could not query remote 'origin': " + err.Error())
	}

	// Use a regular expression to find the line like "HEAD branch: main"
	re := regexp.MustCompile(`HEAD branch:\s*(\S+)`)
	matches := re.FindStringSubmatch(output)

	// matches[0] is the full match ("HEAD branch: main")
	// matches[1] is the captured group ("main")
	if len(matches) < 2 {
		return "", errors.New("could not determine default branch from 'git remote show origin'")
	}

	return matches[1], nil
}
