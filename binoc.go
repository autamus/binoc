package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	fmt.Println("[Parsing Container Blueprints]")
	go repo.Parse(config.Global.Repo.Path, input)

	wg := sync.WaitGroup{}
	wg.Add(1)
	fmt.Println("[Checking Containers for Updates]")
	go update.RunPollWorker(&wg, input, output)

	go func() {
		wg.Wait()
		close(output)
	}()

	for app := range output {
		fmt.Printf("Updating %-30s", app.Data.Name+"...")

		newBranchName := fmt.Sprintf("update-%s", app.Data.Name)
		commitMessage := fmt.Sprintf("Update %s to %s", app.Data.Name, strings.Join(app.Data.LatestVersion.Value, "."))

		_, err := repo.SearchPR(commitMessage, repoOwner, filepath.Base(os.Args[1]), config.Global.Git.Token)
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

		pr, err := repo.SearchPrByBranch(newBranchName, repoOwner, filepath.Base(os.Args[1]), config.Global.Git.Token)
		if err != nil {
			if err.Error() == "not found" {
				err = repo.OpenPR(path, mainBranchName, commitMessage, repoOwner, config.Global.Git.Token, filepath.Base(os.Args[1]))
			}
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = repo.UpdatePR(pr, commitMessage, filepath.Base(os.Args[1]), repoOwner, config.Global.Git.Token)
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
