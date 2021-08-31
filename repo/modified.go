package repo

import (
	"time"

	"github.com/go-git/go-git/v5"
)

func (r *Repo) LastModified(path string) (result time.Time, err error) {
	commits, err := r.backend.Log(&git.LogOptions{
		Order:    git.LogOrderCommitterTime,
		FileName: &path,
	})
	if err != nil {
		return result, err
	}

	commit, err := commits.Next()
	if err != nil {
		return result, err
	}

	commits.Close()
	return commit.Committer.When.UTC(), nil
}
