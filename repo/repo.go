package repo

import (
	"github.com/DataDrake/cuppa/results"
	"github.com/autamus/binoc/parsers"
	"github.com/go-git/go-git/v5"
)

// Result is a reported package and its
// parsed location from the parsed library.
type Result struct {
	Package    parsers.Package
	Parser     parsers.Parser
	LookOutput results.Result
	Path       string
}

type Repo struct {
	Path           string
	enabledParsers map[string]parsers.Parser
	backend        *git.Repository
	gitOptions     *RepoGitOptions
}

type RepoGitOptions struct {
	Name     string
	Username string
	Email    string
	Token    string
}

// Init all enabled parsers from config.
func Init(path string, inputParserNames []string, opts *RepoGitOptions) (result Repo, err error) {
	// Construct enabled parsers for the repository
	result.enabledParsers = make(map[string]parsers.Parser)

	// Loop through input string setting up parsers map.
	for _, parserName := range inputParserNames {
		entry := parsers.AvailableParsers[parserName]
		result.enabledParsers[entry.FileExt] = entry.Parser
	}

	// Open connection to the backend git repository.
	result.gitOptions = opts
	result.backend, err = git.PlainOpen(path)
	result.Path = path
	return result, err
}
