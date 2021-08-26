package repo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Parse parses a single file.
func (r *Repo) Parse(path string) (output Result, err error) {
	match := false
	for ext, parser := range r.enabledParsers {
		match, _ = filepath.Match(ext, filepath.Base(path))
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
	// If package isn't known report unknown.
	if !match {
		return output, errors.New("not a valid package format")
	}
	return output, nil
}

// ParseDir walks through the repository and outputs the parsed values of the spack packages.
func (r *Repo) ParseDir(location string, output chan<- Result) {
	err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		for ext, parser := range r.enabledParsers {
			match, _ := filepath.Match(ext, filepath.Base(path))

                       // If we don't have a match and the parser allows a prefix
                       if !match && parser.AllowsPrefix() {

				fileBasename := filepath.Base(path)
				match = strings.HasPrefix(fileBasename, ext)
                       }

			if match {
				content, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				result, err := parser.Decode(string(content))
				if err != nil {
					fmt.Printf("Parse Error: Couldn't Read --> %s\n", path)
					fmt.Printf("Error: %v\n", err)
					continue
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
