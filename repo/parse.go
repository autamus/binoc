package repo

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Parse parses a single file.
func Parse(path string) (output Result, err error) {
	for ext, parser := range enabledParsers {
		match, _ := filepath.Match(ext, filepath.Base(path))
		if match {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return output, err
			}

			result, err := parser.Decode(string(content))
			if err != nil {
				return output, err
			}
			output = Result{Parser: parser, Package: result, Path: path}
			break
		}
	}

	return output, errors.New("not a valid package format")
}

// ParseDir walks through the repository and outputs the parsed values of the spack packages.
func ParseDir(location string, output chan<- Result) {
	err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		for ext, parser := range enabledParsers {
			match, _ := filepath.Match(ext, filepath.Base(path))
			if match {
				content, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				result, err := parser.Decode(string(content))
				if err != nil {
					return err
				}

				output <- Result{Parser: parser, Package: result, Path: path}
				break
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	close(output)
}
