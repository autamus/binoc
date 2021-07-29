package repo

import (
	"github.com/go-git/go-git/v5"
)

func Reset(path string) (err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	ref, err := r.Head()
	if err != nil {
		return err
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	parent, err := commit.Parents().Next()
	if err != nil {
		return err
	}

	return w.Reset(&git.ResetOptions{
		Commit: parent.Hash,
	})
}
