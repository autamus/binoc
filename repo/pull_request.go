package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

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

// SearchPR searches for the a pull request with the input name in the given repository + owner.
func SearchPR(prTitle, repoOwner, repoName, gitToken string) (state string, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	result, _, err := client.Search.Issues(ctx, fmt.Sprintf("%s type:pr repo:%s/%s", prTitle, repoOwner, repoName), &github.SearchOptions{})
	if err != nil {
		return state, err
	}
	if len(result.Issues) < 1 {
		return "not found", nil
	}
	return *result.Issues[0].State, nil
}
