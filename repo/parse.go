package repo

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/autamus/go-parspack"
)

// Parse walks through the repository and outputs the parsed values of the spack packages.
func Parse(location string, output chan<- Result) {
	err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(filepath.Base(path), ".py") {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			result, err := parspack.Decode(string(content))
			if err != nil {
				return err
			}

			output <- Result{Data: result, Path: path}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	close(output)
}
