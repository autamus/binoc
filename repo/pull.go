package repo

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func Pull(path string, gitUsername string, gitToken string) (err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	branchName, err := GetBranchName(path)
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
				Username: gitUsername,
				Password: gitToken,
			},
		},
	)

	return err
}
