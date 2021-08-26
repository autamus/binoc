package repo

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
)

func (r *Repo) LastModified(path string) (result time.Time, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	commits, err := r.backend.Log(&git.LogOptions{
		Order:    git.LogOrderCommitterTime,
		FileName: &path,
	})
	if err != nil {
		return result, err
	}

	commit, err := commits.Next()
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	commits.Close()
	return commit.Committer.When.UTC(), nil
}
