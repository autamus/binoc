package repo

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func (r *Repo) getURL() (url string, err error) {
	remotes, err := r.backend.Remotes()
	if err != nil {
		return url, err
	}
	return remotes[0].Config().URLs[0], nil
}

func (r *Repo) getOwnerName() (repoOwner, repoName string, err error) {
	url, err := r.getURL()
	if err != nil {
		return repoOwner, repoName, err
	}
	repoName = strings.TrimSuffix(filepath.Base(url), filepath.Ext(url))
	repoOwner = filepath.Base(filepath.Dir(url))
	return repoOwner, repoName, nil
}

// OpenPR opens a pull request from the input branch to the destination branch.
func (r *Repo) OpenPR(mainBranch, prTitle string) (err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: r.gitOptions.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	branchName, err := r.GetBranchName()
	if err != nil {
		return err
	}

	pr := &github.NewPullRequest{
		Title:               github.String(prTitle),
		Body:                github.String(prTitle),
		Head:                github.String(branchName),
		Base:                github.String(mainBranch),
		MaintainerCanModify: github.Bool(true),
	}

	repoOwner, repoName, err := r.getOwnerName()
	if err != nil {
		return err
	}

	_, _, err = client.PullRequests.Create(ctx, repoOwner, repoName, pr)
	return err
}

// UpdatePR updates the Title of a PR.
func (r *Repo) UpdatePR(pr github.Issue, title string) (err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: r.gitOptions.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	newPR := github.PullRequest{
		Title: &title,
		Body:  &title,
	}

	repoOwner, repoName, err := r.getOwnerName()
	if err != nil {
		return err
	}

	_, _, err = client.PullRequests.Edit(ctx, repoOwner, repoName, *pr.Number, &newPR)
	return err
}

// SearchPR searches for the a pull request with the input name in the given repository + owner.
func (r *Repo) SearchPR(prTitle string) (pr github.Issue, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: r.gitOptions.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repoOwner, repoName, err := r.getOwnerName()
	if err != nil {
		return pr, err
	}

	result, _, err := client.Search.Issues(ctx, fmt.Sprintf("%s type:pr repo:%s/%s", prTitle, repoOwner, repoName), &github.SearchOptions{})
	if err != nil {
		return pr, err
	}
	if len(result.Issues) > 0 {
		for i, issue := range result.Issues {
			if *issue.Title == prTitle {
				return result.Issues[i], nil
			}
		}
	}
	return pr, errors.New("not found")
}

// SearchPrByBranch checks to see if there is an existing PR based on a specific branch
// and if so returns the name.
func (r *Repo) SearchPrByBranch(branchName string) (pr github.Issue, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: r.gitOptions.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repoOwner, repoName, err := r.getOwnerName()
	if err != nil {
		return pr, err
	}

	result, _, err := client.Search.Issues(
		ctx,
		fmt.Sprintf("head:%s type:pr repo:%s/%s", branchName, repoOwner, repoName), &github.SearchOptions{})
	if err != nil {
		return pr, err
	}
	if len(result.Issues) > 0 {
		for i, issue := range result.Issues {
			if issue.GetState() == "open" {
				return result.Issues[i], nil
			}
		}
	}
	return pr, errors.New("not found")
}
