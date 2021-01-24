package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alecbcs/lookout/update"
	"github.com/autamus/binoc/config"
	"github.com/autamus/binoc/repo"
	"github.com/autamus/go-parspack/pkg"
)

func main() {
	path := filepath.Join(config.Global.Repos.Path, filepath.Base(os.Args[1]))
	err := repo.Clone(os.Args[1], path)
	if err != nil {
		if err.Error() == "repository already exists" {
			err = repo.Pull(path)
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	pkgs := make(chan pkg.Package)
	go repo.Parse(path, pkgs)

	update.Init(config.Global.GitHub.Token)

	for app := range pkgs {
		fmt.Printf("Parsed: %s\n", app.Name)
		result, found := update.CheckUpdate(app.URL)
		if found {
			if result.Version.Compare(app.LatestVersion.Value) == 0 {
				fmt.Println("Up-To-Date")
			} else {
				fmt.Println("Out-Of-Date")
			}
		} else {
			fmt.Println("Not Found")
		}
	}
}
