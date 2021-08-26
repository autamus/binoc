package repo

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func (r *Repo) Pull() (err error) {
	w, err := r.backend.Worktree()
	if err != nil {
		return err
	}

	branchName, err := r.GetBranchName()
	if err != nil {
		return err
	}

	err = w.Pull(
		&git.PullOptions{
			RemoteName:    "origin",
			SingleBranch:  true,
			Force:         true,
			ReferenceName: plumbing.NewBranchReferenceName(branchName),
			Auth: &http.BasicAuth{
				Username: r.gitOptions.Username,
				Password: r.gitOptions.Token,
			},
		},
	)
	return err
}
