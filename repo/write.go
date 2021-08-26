package repo

import (
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// UpdatePackage patches the package with the current updated package data.
func (r *Repo) UpdatePackage(pkg Result) (err error) {
	file, err := os.Create(pkg.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = pkg.Package.UpdatePackage(pkg.LookOutput)
	if err != nil {
		return err
	}

	result, err := pkg.Parser.Encode(pkg.Package)
	if err != nil {
		return err
	}

	_, err = file.WriteString(result)
	if err != nil {
		return err
	}

	err = file.Sync()
	return err
}

// Commit performs a git commit on the repository.
func (r *Repo) Commit(commitMessage string) (err error) {
	w, err := r.backend.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	commit, err := w.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  r.gitOptions.Name,
			Email: r.gitOptions.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	_, err = r.backend.CommitObject(commit)
	if err != nil {
		return err
	}

	return nil
}
