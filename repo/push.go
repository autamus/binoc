package repo

import (
	"context"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
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

// OpenPR opens a pull request from the input branch to the destination branch.
func OpenPR(path, mainBranch, prTitle, repoOwner, gitToken, repoName string) (err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	h, err := r.Head()
	if err != nil {
		return err
	}

	pr := &github.NewPullRequest{
		Title:               github.String(prTitle),
		Body:                github.String(prTitle),
		Head:                github.String(strings.TrimPrefix(h.Name().String(), "refs/heads/")),
		Base:                github.String(mainBranch),
		MaintainerCanModify: github.Bool(true),
	}

	_, _, err = client.PullRequests.Create(ctx, repoOwner, repoName, pr)
	return err
}
