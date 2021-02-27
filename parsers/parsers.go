package parsers

import (
	"reflect"
	"strings"

	"github.com/alecbcs/cuppa/results"
	"github.com/alecbcs/cuppa/version"
)

// Parser is a universal parser interface implemented by
// all parsers in Binoc
type Parser interface {
	Decode(content string) (pkg Package, err error)
	Encode(pkg Package) (result string, err error)
}

// Package is a universal package interface for working
// with packages in Binoc.
type Package interface {
	AddVersion(results.Result) (err error)
	GetLatestVersion() (result version.Version)
}

type entry struct {
	FileExt string
	Parser  Parser
}

var (
	// AvailableParsers is a map containing all available parsing engines.
	AvailableParsers map[string]entry
)

func registerParser(parser Parser, fileExt string) {
	name := strings.ToLower(reflect.ValueOf(parser).Type().Name())
	AvailableParsers[name] = entry{fileExt: fileExt, parser: parser}
}
