package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/autamus/go-parspack/pkg"

	"github.com/autamus/binoc/repo"
)

func main() {
	err := repo.Clone(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	pkgs := make(chan pkg.Package)
	go repo.Parse(filepath.Join(".binoc/sources/", filepath.Base(os.Args[1])), pkgs)

	for app := range pkgs {
		fmt.Printf("Parsed: %s\n", app.Name)
	}
}
