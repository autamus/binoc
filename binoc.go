package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/autamus/binoc/config"
	"github.com/autamus/binoc/display"
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
	fmt.Println("[Pulling Upstream Repository]")

	path := filepath.Join(config.Global.Repos.Path, filepath.Base(os.Args[1]))
	repoOwner := filepath.Base(filepath.Dir(os.Args[1]))
	err := repo.Clone(os.Args[1], path)
	if err != nil {
		if err.Error() == "repository already exists" {
			err = repo.Pull(path)
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	input := make(chan repo.Result, 20)
	output := make(chan repo.Result, 20)
	fmt.Println("[Parsing Container Blueprints]")
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
		doneChan := make(chan int, 1)
		wg := sync.WaitGroup{}
		wg.Add(1)

		// Display Spinner on Update.
		go display.SpinnerWait(doneChan, "Updating "+app.Data.Name+"...", &wg)

		newBranchName := fmt.Sprintf("update-%s", app.Data.Name)
		commitMessage := fmt.Sprintf("Update %s to %s", app.Data.Name, strings.Join(app.Data.LatestVersion.Value, "."))

		state, err := repo.SearchPR(commitMessage, repoOwner, filepath.Base(os.Args[1]), config.Global.Git.Token)
		if err != nil {
			log.Fatal(err)
		}

		if state != "not found" {
			doneChan <- 0
			wg.Wait()
			fmt.Println()
			continue
		}

		mainBranchName, err := repo.GetBranchName(path)
		if err != nil {
			log.Fatal(err)
		}

		err = repo.CreateBranch(path, newBranchName)
		if err != nil {
			log.Fatal(err)
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

		err = repo.OpenPR(path, mainBranchName, commitMessage, repoOwner, config.Global.Git.Token, filepath.Base(os.Args[1]))
		if err != nil {
			log.Fatal(err)
		}

		err = repo.SwitchBranch(path, mainBranchName)
		if err != nil {
			log.Fatal(err)
		}

		doneChan <- 0
		wg.Wait()
		fmt.Println()
	}

}
