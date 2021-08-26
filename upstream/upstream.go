package upstream

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/autamus/binoc/parsers"
	"github.com/google/go-github/v38/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func GetPackage(upstreamTemplatePath, packageName, token string) (p parsers.Package, modified time.Time, err error) {
	concreteLink := strings.ReplaceAll(upstreamTemplatePath, "{{package}}", packageName)
	resp, err := http.Get(concreteLink)
	if err != nil {
		return nil, modified, err
	}
	if resp.StatusCode != 200 {
		return nil, modified, errors.New("invalid upstream path")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, modified, err
	}
	for _, p := range parsers.AvailableParsers {
		// Attempt to parse file from upstream link
		newPkg, err := p.Parser.Decode(string(body))
		if err != nil {
			continue
		}
		// Get modified date if link is that of a GitHub repository
		if strings.HasPrefix(upstreamTemplatePath, "https://raw.githubusercontent.com/") {
			concreteLink := strings.ReplaceAll(upstreamTemplatePath, "{{package}}", packageName)
			data := strings.Split(strings.TrimPrefix(concreteLink, "https://raw.githubusercontent.com/"), "/")
			owner := data[0]
			repo := data[1]
			ctx := context.Background()
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(ctx, ts)

			client := github.NewClient(tc)
			commits, _, err := client.Repositories.ListCommits(ctx, owner, repo, &github.CommitsListOptions{
				Path: strings.Join(data[3:], "/"),
				ListOptions: github.ListOptions{
					Page:    1,
					PerPage: 1,
				},
			})
			if err != nil {
				return newPkg, modified, err
			}
			return newPkg, *commits[0].Commit.Committer.Date, nil
		}
		return newPkg, modified, nil
	}
	return nil, modified, errors.New("unable to parse package")
}
