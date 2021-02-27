package repo

import "github.com/autamus/binoc/parsers"

// Result is a reported package and its
// parsed location from the parsed library.
type Result struct {
	Package parsers.Package
	Parser  parsers.Parser
	Path    string
}

var (
	enabledParsers map[string]parsers.Parser
)

// Init all enabled parsers from config.
func Init(inputParserNames []string) {
	for _, parserName := range inputParserNames {
		entry := parsers.AvailableParsers[parserName]
		enabledParsers[entry.FileExt] = entry.Parser
	}
}
