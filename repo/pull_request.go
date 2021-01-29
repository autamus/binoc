package repo

import (
	"context"
	"errors"
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

// UpdatePR updates the Title of a PR.
func UpdatePR(pr github.Issue, title, repo, repoOwner string, gitToken string) (err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	newPR := github.PullRequest{
		Title: &title,
		Body:  &title,
	}

	_, _, err = client.PullRequests.Edit(ctx, repoOwner, repo, *pr.Number, &newPR)
	return err
}

// SearchPR searches for the a pull request with the input name in the given repository + owner.
func SearchPR(prTitle, repoOwner, repoName, gitToken string) (pr github.Issue, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	result, _, err := client.Search.Issues(ctx, fmt.Sprintf("%s type:pr repo:%s/%s", prTitle, repoOwner, repoName), &github.SearchOptions{})
	if err != nil {
		return pr, err
	}
	if len(result.Issues) < 1 {
		return pr, errors.New("not found")
	}
	return result.Issues[0], nil
}

// SearchPrByBranch checks to see if there is an existing PR based on a specific branch
// and if so returns the name.
func SearchPrByBranch(branchName, repoOwner, repoName, gitToken string) (pr github.Issue, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	result, _, err := client.Search.Issues(
		ctx,
		fmt.Sprintf("head:%s type:pr repo:%s/%s", branchName, repoOwner, repoName), &github.SearchOptions{})
	if err != nil {
		return pr, err
	}
	if len(result.Issues) < 1 {
		return pr, errors.New("not found")
	}
	return result.Issues[0], nil
}
