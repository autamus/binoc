package repo

import (
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// CreateBranch creates a new branch within the input
// repository.
func CreateBranch(path string, name string) (err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	h, err := r.Head()
	if err != nil {
		return err
	}

	ref := plumbing.NewHashReference(plumbing.NewBranchReferenceName(name), h.Hash())
	err = r.Storer.SetReference(ref)

	return err

}

// SwitchBranch switches from the current branch to the
// one with the name provided.
func SwitchBranch(path string, branchName string) (err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	branchRef := plumbing.NewBranchReferenceName(branchName)
	opts := &git.CheckoutOptions{Branch: branchRef}

	err = w.Checkout(opts)
	return err
}

// GetBranchName returns the name of the current branch.
func GetBranchName(path string) (name string, err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return name, err
	}
	h, err := r.Head()
	if err != nil {
		return name, err
	}
	name = strings.TrimPrefix(h.Name().String(), "refs/heads/")

	return name, nil
}
