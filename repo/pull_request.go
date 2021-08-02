package repo

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func getURL(path string) (url string, err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return url, err
	}
	remotes, err := r.Remotes()
	if err != nil {
		return url, err
	}
	return remotes[0].Config().URLs[0], nil
}

func getOwnerName(path string) (repoOwner, repoName string, err error) {
	url, err := getURL(path)
	if err != nil {
		return repoOwner, repoName, err
	}
	repoName = strings.TrimSuffix(filepath.Base(url), filepath.Ext(url))
	repoOwner = filepath.Base(filepath.Dir(url))
	return repoOwner, repoName, nil
}

// OpenPR opens a pull request from the input branch to the destination branch.
func OpenPR(path, mainBranch, prTitle, gitToken string) (err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	branchName, err := GetBranchName(path)
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

	repoOwner, repoName, err := getOwnerName(path)
	if err != nil {
		return err
	}

	_, _, err = client.PullRequests.Create(ctx, repoOwner, repoName, pr)
	return err
}

// UpdatePR updates the Title of a PR.
func UpdatePR(pr github.Issue, path, title, gitToken string) (err error) {
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

	repoOwner, repoName, err := getOwnerName(path)
	if err != nil {
		return err
	}

	_, _, err = client.PullRequests.Edit(ctx, repoOwner, repoName, *pr.Number, &newPR)
	return err
}

// SearchPR searches for the a pull request with the input name in the given repository + owner.
func SearchPR(path, prTitle, gitToken string) (pr github.Issue, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repoOwner, repoName, err := getOwnerName(path)
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
func SearchPrByBranch(path, branchName, gitToken string) (pr github.Issue, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repoOwner, repoName, err := getOwnerName(path)
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
				pullRequest, _, err := client.PullRequests.Get(ctx, repoOwner, repoName, *issue.Number)
				if err != nil {
					return pr, err
				}
				if pullRequest.GetHead().GetLabel() == branchName {
					return result.Issues[i], nil
				}
			}
		}
	}
	return pr, errors.New("not found")
}
