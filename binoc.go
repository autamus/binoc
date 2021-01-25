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
	go repo.Parse(path, input)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go update.RunPollWorker(&wg, input, output)

	go func() {
		wg.Wait()
		close(output)
	}()

	for app := range output {
		newBranchName := fmt.Sprintf("update-%s", app.Data.Name)
		commitMessage := fmt.Sprintf("Update %s to %s", app.Data.Name, strings.Join(app.Data.LatestVersion.Value, "."))
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

		fmt.Println(app.Data.Name)
	}

}
