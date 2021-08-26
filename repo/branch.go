package repo

import (
	"errors"
	"log"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

// CreateBranch creates a new branch within the input
// repository.
func (r *Repo) CreateBranch(branchName string) (err error) {
	h, err := r.backend.Head()
	if err != nil {
		return err
	}

	ref := plumbing.NewHashReference(plumbing.NewBranchReferenceName(branchName), h.Hash())
	err = r.backend.Storer.SetReference(ref)

	return err

}

// SwitchBranch switches from the current branch to the
// one with the name provided.
func (r *Repo) SwitchBranch(branchName string) (err error) {
	w, err := r.backend.Worktree()
	if err != nil {
		return err
	}

	branchRef := plumbing.NewBranchReferenceName(branchName)
	opts := &git.CheckoutOptions{Branch: branchRef}

	err = w.Checkout(opts)
	return err
}

// GetBranchName returns the name of the current branch.
func (r *Repo) GetBranchName() (name string, err error) {
	h, err := r.backend.Head()
	if err != nil {
		return name, err
	}
	name = strings.TrimPrefix(h.Name().String(), "refs/heads/")

	return name, nil
}

// PullBranch attempts to pull the branch from the git origin.
func (r *Repo) PullBranch(branchName string) (err error) {
	localBranchReferenceName := plumbing.NewBranchReferenceName(branchName)
	remoteReferenceName := plumbing.NewRemoteReferenceName("origin", branchName)

	rem, err := r.backend.Remote("origin")
	if err != nil {
		return err
	}

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	found := false
	for _, ref := range refs {
		if ref.Name().IsBranch() && ref.Name() == localBranchReferenceName {
			found = true
		}
	}

	if !found {
		return errors.New("branch not found")
	}

	err = r.backend.CreateBranch(&config.Branch{Name: branchName, Remote: "origin", Merge: localBranchReferenceName})
	if err != nil {
		return err
	}
	newReference := plumbing.NewSymbolicReference(localBranchReferenceName, remoteReferenceName)
	err = r.backend.Storer.SetReference(newReference)
	return err
}
