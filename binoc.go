package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/autamus/binoc/config"
	"github.com/autamus/binoc/repo"
	"github.com/autamus/binoc/update"
)

func main() {
	update.Init(config.Global.Git.Token)
	fmt.Print(` ____  _                  
| __ )(_)_ __   ___   ___ 
|  _ \| | '_ \ / _ \ / __|
| |_) | | | | | (_) | (__ 
|____/|_|_| |_|\___/ \___|
`)
	fmt.Printf("Application Version: v%s\n", config.Global.General.Version)
	fmt.Println()

	input := make(chan repo.Result, 20)
	output := make(chan repo.Result, 20)
	path := config.Global.Repo.Path
	fmt.Println("[Parsing Container Blueprints]")

	// Parse Config Value into list of parser names
	repo.Init(strings.Split(config.Global.Parsers.Loaded, ","))

	// Begin parsing the repository matching file extentions to parsers.
	go repo.Parse(path, input)

	wg := sync.WaitGroup{}
	wg.Add(1)
	fmt.Println("[Checking Containers for Updates]")
	go update.RunPollWorker(&wg, input, output)

	go func() {
		wg.Wait()
		close(output)
	}()

	for app := range output {
		name := app.Package.GetName()

		fmt.Printf("Updating %-30s", name+"...")

		newBranchName := fmt.Sprintf("update-%s", name)
		commitMessage := fmt.Sprintf("Update %s to %s", name, strings.Join(app.Package.GetLatestVersion(), "."))

		_, err := repo.SearchPR(path, commitMessage, config.Global.Git.Token)
		if err == nil {
			fmt.Println("Skipped")
			continue
		}
		if err.Error() != "not found" {
			log.Fatal(err)
		}

		mainBranchName, err := repo.GetBranchName(path)
		if err != nil {
			log.Fatal(err)
		}

		err = repo.PullBranch(path, newBranchName)
		if err != nil {
			if err.Error() == "branch not found" {
				err = repo.CreateBranch(path, newBranchName)
			}
			if err != nil {
				log.Fatal(err)
			}
		}

		err = repo.SwitchBranch(path, newBranchName)
		if err != nil {
			log.Fatal(err)
		}

		err = repo.UpdatePackage(app)
		if err != nil {
			log.Fatal(err)
		}

		err = repo.Commit(path, commitMessage, config.Global.Git.Name, config.Global.Git.Email)
		if err != nil {
			log.Fatal(err)
		}

		err = repo.Push(path, config.Global.Git.Username, config.Global.Git.Token)
		if err != nil {
			log.Fatal(err)
		}

		pr, err := repo.SearchPrByBranch(path, newBranchName, config.Global.Git.Token)
		if err != nil {
			if err.Error() == "not found" {
				err = repo.OpenPR(path, mainBranchName, commitMessage, config.Global.Git.Token)
			}
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = repo.UpdatePR(pr, path, commitMessage, config.Global.Git.Token)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = repo.SwitchBranch(path, mainBranchName)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Done")
	}

}
