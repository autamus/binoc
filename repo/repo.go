package repo

import "github.com/autamus/binoc/parsers"

// Result is a reported package and its
// parsed location from the parsed library.
type Result struct {
	Package parsers.Package
	Parser  parsers.Parser
	Path    string
}
