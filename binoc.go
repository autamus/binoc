package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/autamus/binoc/config"
	"github.com/autamus/binoc/parsers"
	"github.com/autamus/binoc/repo"
	"github.com/autamus/binoc/update"
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/github"
)

func main() {
	update.Init(config.Global.Git.Token)
	fmt.Println()
	fmt.Print(` ____  _                  
| __ )(_)_ __   ___   ___ 
|  _ \| | '_ \ / _ \ / __|
| |_) | | | | | (_) | (__ 
|____/|_|_| |_|\___/ \___|
`)
	fmt.Printf("Application Version: %s\n", config.Global.General.Version)
	fmt.Println()

	input := make(chan repo.Result, 20)
	output := make(chan repo.Result, 20)
	relay := make(chan repo.Result, 20)

	parsed := 0
	updated := 0
	skipped := 0

	path := config.Global.Repo.Path
	if config.Global.General.Action == "true" {
		path = filepath.Join("/github/workspace", path)
	}

	fmt.Println("[Parsing Container Blueprints]")

	// Parse Config Value into list of parser names
	repo, err := repo.Init(path,
		strings.Split(config.Global.Parsers.Loaded, ","),
		&repo.RepoGitOptions{
			Name:     config.Global.Git.Name,
			Username: config.Global.Git.Username,
			Email:    config.Global.Git.Email,
			Token:    config.Global.Git.Token,
		},
	)
	if err != nil {
		printError(err)
	}

	// Set Config Value to Spack Parser
	parsers.SpackUpstreamLink = config.Global.Repo.SpackUpstreamLink

	// Pull Git Repository Updates
	err = repo.Pull()
	if err != nil && err != git.NoErrAlreadyUpToDate {
		printError(err)
	}

	// Begin parsing the repository matching file extentions to parsers.
	go repo.ParseDir(path, relay)
	go func() {
		for app := range relay {
			parsed++
			input <- app
		}
		close(input)
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU()*2; i++ {
		go update.RunPollWorker(
			&wg,
			&repo,
			config.Global.Repo.SpackUpstreamLink,
			config.Global.Git.Token,
			input,
			output,
		)
		wg.Add(1)
	}

	fmt.Println("[Checking Containers for Updates]")

	go func() {
		wg.Wait()
		close(output)
	}()

	// Store the name of the "main" branch that we
	// started on.
	mainBranchName, err := repo.GetBranchName()
	if err != nil {
		printError(err)
	}

	for app := range output {
		name := app.Package.GetName()
		if strings.HasPrefix(app.LookOutput.Name, "spack") {
			fmt.Printf("Fixing %-32s", name+"...")
		} else {
			fmt.Printf("Updating %-30s", name+"...")
		}

		var commitMessage, newBranchName string
		var pr github.Issue

		// Only run git checkouts, commits, if binoc is managing PRs
		if config.Global.PR.Skip == "false" {
			if strings.HasPrefix(app.LookOutput.Name, "spack") {
				newBranchName = fmt.Sprintf("%sfix-%s", config.Global.Branch.Prefix, name)
				commitMessage = fmt.Sprintf("Pull Changes to %s from Upstream", name)
			} else {
				newBranchName = fmt.Sprintf("%supdate-%s", config.Global.Branch.Prefix, name)
				commitMessage = fmt.Sprintf("Update %s to %s", name, strings.Join(app.LookOutput.Version, "."))
			}

			// Search for previous open pull requests so that we don't create duplicates.
			pr, err = repo.SearchPR(commitMessage)
			if err != nil && err.Error() != "not found" {
				printError(err)
			}
			if err == nil {
				blacklistFound := false
				for _, label := range pr.Labels {
					if *label.Name == config.Global.PR.IgnoreLabel {
						blacklistFound = true
					}
				}
				if *pr.State == "open" || blacklistFound {
					fmt.Println("Skipped")
					skipped++
					continue
				}
			}

			// Pull an existing branch to update if possible.
			err = repo.PullBranch(newBranchName)
			if err != nil {
				if err.Error() == "branch not found" {
					err = repo.CreateBranch(newBranchName)
				}
				if err != nil {
					printError(err)
				}
			}

			err = repo.SwitchBranch(newBranchName)
			if err != nil {
				printError(err)
			}
		}

		// Updating the package is run regardless of pr_skip
		err = repo.UpdatePackage(app)
		if err != nil {
			printError(err)
		}

		// If we are not managing prs, continue in loop to update
		if config.Global.PR.Skip == "false" {

			err = repo.Commit(commitMessage)
			if err != nil {
				printError(err)
			}

			err = repo.Push()
			if err != nil {
				if err != nil {
					printError(err)
				}
			}

			pr, err = repo.SearchPrByBranch(newBranchName)
			if err == nil && *pr.State == "open" {
				err = repo.UpdatePR(pr, commitMessage)
				if err != nil {
					printError(err)
				}
			} else {
				if err != nil && err.Error() != "not found" {
					printError(err)
				}
				err = repo.OpenPR(mainBranchName, commitMessage)
				if err != nil {
					printError(err)
				}
			}
			err = repo.SwitchBranch(mainBranchName)
			if err != nil {
				printError(err)
			}

		}

		fmt.Println("Done")
		updated++

		// We only need to sleep if we are submitting PRs
		if config.Global.PR.Skip == "false" {
			time.Sleep(5 * time.Second)
		}
	}
	fmt.Println()
	fmt.Println("[Scan Results]")
	fmt.Printf("%-5d Packages Parsed\n", parsed)
	fmt.Printf("%-5d Packages Updated\n", updated)
	fmt.Printf("%-5d Packages Skipped\n", skipped)
	fmt.Println()
}

func printError(err error) {
	fmt.Println("Error")
	log.Fatal(err)
}
