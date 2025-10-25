# git-smart ⚡️

[![Go Version][go-shield]][]
[![Release Version][release-shield]][]

**Tired of the repetitive `git checkout main`, `git pull`, `git checkout feature`, `git rebase main` dance?** `git-smart sync` automates syncing your feature branch with the default branch in one command.

---

## ✨ Demo



---

## 🤔 The Problem

Keeping your feature branches up-to-date with the latest changes from `main` (or `master`) is essential, but it's a tedious, multi-step process:

1.  `git stash` (if you have local changes)
2.  `git checkout main`
3.  `git pull`
4.  `git checkout your-feature-branch`
5.  `git rebase main`
6.  `git stash pop` (if you stashed)
7.  Resolve any conflicts...

It's easy to forget a step or make a mistake, especially when you're doing it multiple times a day.

---

## 🎉 The Solution: `git-smart sync`

`git-smart` provides a single, intelligent command:

```bash
git-smart sync
````

It automatically handles the entire workflow for you:

  * Detects your default branch (`main` or `master`).
  * Stashes uncommitted changes safely.
  * Checks out and pulls the default branch.
  * Checks back out your original branch.
  * Rebases the default branch onto yours.
  * Applies your stashed changes.
  * Clearly reports conflicts if they occur.

Spend less time juggling Git commands and more time coding\!

-----

## 🚀 Installation

### macOS / Linux (Homebrew)

```bash
# First, tap the repository (only need to do this once):
brew tap AEGON247/tap

# Then, install git-smart:
brew install git-smart
```

### Windows (Scoop)

*(Coming Soon - Requires GoReleaser configuration)*

### Manual Download

Download the latest binary for your operating system from the [**GitHub Releases**][release-url] page and add it to your system's PATH.

-----

## ⌨️ Usage

Simply navigate to any directory managed by Git and run:

```bash
git-smart sync
```

The tool will guide you through the process with clear, colored output. If a rebase conflict occurs, it will stop and instruct you on how to resolve it.

-----

## 🤝 Contributing

Contributions are welcome\! Please feel free to submit issues or pull requests on the [GitHub repository][https://github.com/AEGON247/git-smart].

-----


[release-url]: https://www.google.com/search?q=https://github.com/AEGON247/git-smart/releases
[repo-url]: https://www.google.com/search?q=https://github.com/AEGON247/git-smart
[license-url]: https://www.google.com/search?q=https://github.com/AEGON247/git-smart/blob/main/LICENSE