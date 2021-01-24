package repo

import (
	"github.com/go-git/go-git/v5"
)

// Clone performs a git clone on the upstream url repository.
func Clone(url string, path string) (err error) {
	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	return err
}

// Pull performs a git pull on the upstream url repository.
// If the repo is already up to date the function exits
// and does not report an error.
func Pull(path string) (err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err.Error() != "already up-to-date" {
		return err
	}

	return nil
}
