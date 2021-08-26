package repo

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Push performs a "git push" on the repository.
func Push(path string, gitUsername string, gitToken string) (err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	h, err := r.Head()
	if err != nil {
		return err
	}
	// Generate <src>:<dest> reference string
	refStr := h.Name().String() + ":" + h.Name().String()
	// Push Branch to Origin
	err = r.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{config.RefSpec(refStr)},
		Auth: &http.BasicAuth{
			Username: gitUsername,
			Password: gitToken,
		},
	})
	return err
}
