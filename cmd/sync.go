// cmd/sync.go
package cmd

import (
	"fmt"
	"os" // Make sure 'os' is imported
	"regexp"

	"github.com/AEGON247/git-smart/internal/git" // Adjust path if needed
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Stashes changes, syncs with the default branch (main/master), and pops stash.",
	Long: `Automates the common workflow of keeping a feature branch up-to-date.
It performs the following steps:
1. Gets current branch name.
2. Stashes uncommitted changes (if any).
3. Checks out the default branch (e.g., main or master).
4. Pulls the latest changes for the default branch.
5. Checks out the original feature branch.
6. Rebases the default branch onto the feature branch.
7. Pops the stash (if anything was stashed).`,
	Run: func(cmd *cobra.Command, args []string) {
		// --- 0. Setup ---
		stepColor := color.New(color.FgCyan).Add(color.Bold)
		successColor := color.New(color.FgGreen)
		errorColor := color.New(color.FgRed)
		infoColor := color.New(color.FgYellow)

		// Check if it's a git repo first
		if !git.IsRepo() {
			errorColor.Println("Error: This is not a git repository.")
			os.Exit(1)
		}

		// --- Get Default Branch Name ---
		stepColor.Println("STEP 0: Determining default branch name...")
		defaultBranch, err := git.GetDefaultBranch()
		if err != nil {
			errorColor.Println("Error getting default branch:", err)
			os.Exit(1)
		}
		successColor.Println("Default branch detected:", defaultBranch)

		// --- 1. Get Current Branch ---
		stepColor.Println("\nSTEP 1: Checking current branch...")
		currentBranch, err := git.GetCurrentBranch()
		if err != nil {
			errorColor.Println("Error getting current branch:", err)
			os.Exit(1)
		}

		// --- Handle edge case: User is on the default branch ---
		if currentBranch == defaultBranch {
			infoColor.Printf("You are already on the '%s' branch. Pulling latest changes...\n", defaultBranch)
			output, err := git.Run("pull")
			if err != nil {
				errorColor.Printf("Error pulling '%s': %s\n", defaultBranch, err)
				os.Exit(1)
			}
			successColor.Printf("Successfully pulled '%s'. You are up to date.\n", defaultBranch)
			fmt.Println(output)
			return // We are done!
		}

		successColor.Println("On feature branch:", currentBranch)

		// --- 2. Stash Changes (if any) ---
		stepColor.Println("\nSTEP 2: Checking for uncommitted changes...")
		status, err := git.Run("status", "--porcelain")
		if err != nil {
			errorColor.Println("Error checking status:", err)
			os.Exit(1)
		}

		didStash := false
		if status != "" {
			infoColor.Println("Uncommitted changes found. Stashing...")
			_, err := git.Run("stash")
			if err != nil {
				errorColor.Println("Error stashing changes:", err)
				os.Exit(1)
			}
			didStash = true
			successColor.Println("Changes stashed.")
		} else {
			successColor.Println("No changes to stash.")
		}

		// --- 3. Checkout Default Branch ---
		stepColor.Printf("\nSTEP 3: Checking out '%s'...\n", defaultBranch)
		if _, err := git.Run("checkout", defaultBranch); err != nil {
			errorColor.Printf("Error checking out '%s': %s\n", defaultBranch, err)
			// Attempt to undo stash if we stashed
			if didStash {
				infoColor.Println("Attempting to restore stash...")
				git.Run("stash", "pop")
			}
			os.Exit(1)
		}

		// --- 4. Pull Default Branch ---
		stepColor.Printf("\nSTEP 4: Pulling latest changes for '%s'...\n", defaultBranch)
		if _, err := git.Run("pull"); err != nil {
			errorColor.Printf("Error pulling '%s': %s\n", defaultBranch, err)
			infoColor.Println("Attempting to switch back to your branch...")
			git.Run("checkout", currentBranch) // Try to go back
			if didStash {
				infoColor.Println("Attempting to restore stash...")
				git.Run("stash", "pop")
			}
			os.Exit(1)
		}
		successColor.Printf("Pulled '%s' successfully.\n", defaultBranch)

		// --- 5. Checkout Feature Branch ---
		stepColor.Println("\nSTEP 5: Checking out feature branch...")
		if _, err := git.Run("checkout", currentBranch); err != nil {
			errorColor.Println("Error checking out feature branch:", err)
			if didStash {
				infoColor.Println("Attempting to restore stash on default branch...")
				git.Run("stash", "pop") // Just pop, we're on default
			}
			os.Exit(1)
		}

		// --- 6. Rebase ---
		stepColor.Printf("\nSTEP 6: Rebasing '%s' onto your branch...\n", defaultBranch)
		if output, err := git.Run("rebase", defaultBranch); err != nil {
			// THIS IS A REBASE CONFLICT
			errorColor.Println("REBASE FAILED: You have conflicts.")
			infoColor.Println("--- Git Output ---")
			infoColor.Println(output)
			infoColor.Println("-----------------")
			infoColor.Println("Please fix the conflicts and then run 'git rebase --continue'.")
			if didStash {
				infoColor.Println("Your stashed changes were NOT applied. Run 'git stash pop' after your rebase is complete.")
			}
			os.Exit(1)
		}
		successColor.Println("Rebase successful.")

		// --- 7. Pop Stash ---
		if didStash {
			stepColor.Println("\nSTEP 7: Applying stashed changes...")
			popOutput, err := git.Run("stash", "pop")
			if err != nil {
				// Check if the error is just "No stash found" - might happen if stash was empty or applied during conflict resolution
				if regexp.MustCompile(`(?i)No stash found|No stash entries found`).MatchString(popOutput) || regexp.MustCompile(`(?i)Did not need to pop stash`).MatchString(popOutput) {
					successColor.Println("Stash was empty or already applied.")
				} else {
					errorColor.Println("Error popping stash. Your stash is still saved.")
					infoColor.Println(popOutput) // Print the error message from git
					os.Exit(1)
				}
			} else {
				successColor.Println("Stashed changes applied.")
			}
		}

		// --- Done! ---
		successColor.Printf("\n✅ All done! Your branch is synced with '%s'.\n", defaultBranch)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// You can add flags here later if needed, e.g., --merge instead of --rebase
	// syncCmd.Flags().BoolP("merge", "m", false, "Use merge instead of rebase")
}
