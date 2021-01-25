package repo

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Push performs a "git push" on the repository.
func Push(path string, gitUsername string, gitToken string) (err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: gitUsername,
			Password: gitToken,
		},
	})
	return err
}
