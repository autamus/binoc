package repo

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Push performs a "git push" on the repository.
func (r *Repo) Push() (err error) {
	h, err := r.backend.Head()
	if err != nil {
		return err
	}
	// Generate <src>:<dest> reference string
	refStr := h.Name().String() + ":" + h.Name().String()
	// Push Branch to Origin
	err = r.backend.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{config.RefSpec(refStr)},
		Auth: &http.BasicAuth{
			Username: r.gitOptions.Username,
			Password: r.gitOptions.Token,
		},
	})
	return err
}
