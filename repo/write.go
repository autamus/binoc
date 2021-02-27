package repo

import (
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// UpdatePackage patches the package with the current updated package data.
func UpdatePackage(pkg Result) (err error) {
	err = pkg.Package.AddVersion(pkg.LookOutput)
	if err != nil {
		return err
	}

	file, err := os.Create(pkg.Path)
	if err != nil {
		return err
	}
	defer file.Close()

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
func Commit(path string, commitMessage string, gitName string, gitEmail string) (err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	commit, err := w.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitName,
			Email: gitEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	_, err = r.CommitObject(commit)
	if err != nil {
		return err
	}

	return nil
}
