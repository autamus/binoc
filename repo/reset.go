package repo

import (
	"github.com/go-git/go-git/v5"
)

func (r *Repo) Reset() (err error) {
	w, err := r.backend.Worktree()
	if err != nil {
		return err
	}

	ref, err := r.backend.Head()
	if err != nil {
		return err
	}
	commit, err := r.backend.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	parent, err := commit.Parents().Next()
	if err != nil {
		return err
	}

	err = w.Reset(&git.ResetOptions{
		Commit: parent.Hash,
	})
	if err != nil {
		return err
	}
	_, err = w.Add(".")
	return err
}
